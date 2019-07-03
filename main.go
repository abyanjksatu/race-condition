package main

import (
	"fmt"
	"sync"
)

type safeNumber struct {
	val int
	m   sync.Mutex
}

func main() {
	fmt.Println("Blocking With waitgroups")
	// The most straightforward way of solving a data race, is to
	// block read access until the write operation has been completed
	fmt.Println(blockingWithWaitgroups())

	fmt.Println("Blocking With channels")
	// Blocking inside the getNumber function, although simple, would get
	// troublesome if we want to call the function repeatedly. The next
	// method follows a more flexible approach towards blocking.
	fmt.Println(blockingWithChannel())

	fmt.Println("Returning a channels")
	// The code is blocked until something gets pushed into the returned channel
	// As opposed to the previous method, we block in the main function, instead
	// of the function itself
	i := <-returningWithChannel()
	fmt.Println(i)

	fmt.Println("Using Mutex")
	// Until now, we had decided that the value of i should only be read after
	// the write operation has finished. Let’s now think about the case, where
	// we don’t care about the order of reads and writes, we only require that
	// they do not occur simultaneously. If this sounds like your use case,
	// then you should consider using a mutex
	fmt.Println(useMutex())
}

func blockingWithWaitgroups() int {
	var i int
	// Initialize a waitgroup variable
	var wg sync.WaitGroup
	// `Add(1) signifies that there is 1 task that we need to wait for
	wg.Add(1)
	go func() {
		i = 5
		// Calling `wg.Done` indicates that we are done with the task we are waiting fo
		wg.Done()
	}()
	// `wg.Wait` blocks until `wg.Done` is called the same number of times
	// as the amount of tasks we have (in this case, 1 time)
	wg.Wait()
	return i
}

func blockingWithChannel() int {
	var i int
	// Create a channel to push an empty struct to once we're done
	done := make(chan struct{})
	go func() {
		i = 5
		// Push an empty struct once we're done
		done <- struct{}{}
	}()
	// This statement blocks until something gets pushed into the `done` channel
	<-done
	return i
}

func returningWithChannel() <-chan int {
	// create the channel
	c := make(chan int)
	go func() {
		// push the result into the channel
		c <- 5
	}()
	// immediately return the channel
	return c
}

func (i *safeNumber) get() int {
	// The `Lock` method of the mutex blocks if it is already locked
	// if not, then it blocks other calls until the `Unlock` method is called
	i.m.Lock()
	// Defer `Unlock` until this method returns
	defer i.m.Unlock()
	// Return the value
	return i.val
}

func (i *safeNumber) set(val int) {
	// Similar to the `Get` method, except we Lock until we are done
	// writing to `i.val`
	i.m.Lock()
	defer i.m.Unlock()
	i.val = val
}

func useMutex() int {
	// Create an instance of `safeNumber`
	i := &safeNumber{}
	// Use `Set` and `Get` instead of regular assignments and reads
	// We can now be sure that we can read only if the write has completed, or vice versa
	go func() {
		i.set(5)
	}()
	return i.get()
}
