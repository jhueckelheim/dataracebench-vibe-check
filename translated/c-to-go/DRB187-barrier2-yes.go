/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

not a barrier. either thread can exit before the other thread
enters. So race on x can occur.
Data race pair: x@39:7:W vs. x@51:7:W
*/

package main

import (
	"fmt"
	"sync"
)

var (
	l0, l1 sync.Mutex
	x      = 1
)

func main() {
	// Initialize locks
	l0.Lock()
	l1.Lock()

	var wg sync.WaitGroup
	wg.Add(2)

	// Thread 0
	go func() {
		defer wg.Done()
		tid := 0

		// Initial barrier equivalent

		if tid == 0 {
			x = 0
		}

		// Faulty synchronization - each thread locks/unlocks its own lock
		// This provides no cross-thread synchronization!
		if tid == 0 {
			l0.Unlock() // Release own lock
			l0.Lock()   // Immediately re-acquire own lock - no waiting!
		}

		// Race: Thread 1 might write to x concurrently
	}()

	// Thread 1
	go func() {
		defer wg.Done()
		tid := 1

		// Initial barrier equivalent

		// Faulty synchronization - each thread locks/unlocks its own lock
		if tid == 1 {
			l1.Unlock() // Release own lock
			l1.Lock()   // Immediately re-acquire own lock - no waiting!
		}

		if tid == 1 {
			x = 1 // Race: Thread 0 might write to x concurrently
		}

		// Cleanup
		l1.Unlock()
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
