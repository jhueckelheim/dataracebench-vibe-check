/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

Two threads sync using one lock. No race.
*/

package main

import (
	"fmt"
	"sync"
)

var (
	l sync.Mutex
	x = 1
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Thread 0
	go func() {
		defer wg.Done()
		tid := 0

		if tid == 0 {
			l.Lock() // Thread 0 acquires the lock
		}

		// Barrier equivalent

		if tid == 0 {
			x = 0
			l.Unlock() // Thread 0 releases the lock
		}
	}()

	// Thread 1
	go func() {
		defer wg.Done()
		tid := 1

		// Barrier equivalent

		if tid == 1 {
			l.Lock()   // Thread 1 waits for thread 0 to release the lock
			l.Unlock() // Thread 1 immediately releases it
			x = 1
		}
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
