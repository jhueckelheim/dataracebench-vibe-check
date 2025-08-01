/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

sync with busy wait loop using atomic. No data race pair.
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var x, s int32
	x = 0
	s = 0

	var wg sync.WaitGroup
	wg.Add(2)

	// Section 1
	go func() {
		defer wg.Done()
		x = 1
		// Atomic write with sequential consistency
		atomic.StoreInt32(&s, 1)
	}()

	// Section 2
	go func() {
		defer wg.Done()
		done := int32(0)
		for done == 0 {
			// Atomic read with sequential consistency
			done = atomic.LoadInt32(&s)
		}
		x = 2
	}()

	wg.Wait()
	fmt.Printf("%d\n", x)
}
