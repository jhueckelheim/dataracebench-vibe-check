/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

uses locks to create a barrier that is used once.
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
	var wg sync.WaitGroup
	wg.Add(2)

	// Initialize locks - both start locked
	l0.Lock()
	l1.Lock()

	// Thread 0
	go func() {
		defer wg.Done()

		// Initial barrier - thread 0 already has l0
		x = 0

		// Synchronization sequence from original C
		l0.Unlock() // Release own lock
		l1.Lock()   // Acquire thread 1's lock - wait for its signal
	}()

	// Thread 1
	go func() {
		defer wg.Done()

		// Synchronization sequence from original C
		l0.Lock()   // Acquire thread 0's lock - wait for its signal
		l1.Unlock() // Release own lock

		x = 1
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
