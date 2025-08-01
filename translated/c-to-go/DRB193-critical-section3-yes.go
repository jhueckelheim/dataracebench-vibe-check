/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

Race due to different critical section names
Data race pair: x@27:7:W vs. x@44:7:W
Data race pair: s@30:9:W vs. s@40:15:R
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var x, s int
	var mutexA, mutexB sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	// Section 1 - uses critical section A
	go func() {
		defer wg.Done()
		x = 1

		mutexA.Lock() // Critical section A
		s = 1
		mutexA.Unlock()
	}()

	// Section 2 - uses critical section B
	go func() {
		defer wg.Done()
		done := 0

		for done == 0 {
			mutexB.Lock() // Critical section B - different mutex!
			if s != 0 {
				done = 1
			}
			mutexB.Unlock()
		}
		x = 2 // Race: both threads write to x using different mutexes for signaling
	}()

	wg.Wait()
	fmt.Printf("%d\n", x)
}
