/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

Race because the write to s is not protected by atomic
Data race pair: s@26:7:W vs. s@34:16:R
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var x int32
	var s int32
	x = 0
	s = 0

	var wg sync.WaitGroup
	wg.Add(2)

	// Section 1
	go func() {
		defer wg.Done()
		x = 1
		// Race: Non-atomic write to s
		s = 1
	}()

	// Section 2
	go func() {
		defer wg.Done()
		done := int32(0)
		for done == 0 {
			// Race: Atomic read of s while another goroutine does non-atomic write
			done = atomic.LoadInt32(&s)
		}
		x = 2
	}()

	wg.Wait()
	fmt.Printf("%d\n", x)
}
