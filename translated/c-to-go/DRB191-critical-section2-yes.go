/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

Race due to different critical section names
Data race pair: size@34:11:W vs. size@49:11:W
*/

package main

import (
	"fmt"
	"sync"
)

var (
	cap    = 10
	size   = 0
	mutexA sync.Mutex // Critical section A
	mutexB sync.Mutex // Critical section B
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Producer section - uses critical section A
	go func() {
		defer wg.Done()
		r := uint(0)

		for {
			mutexA.Lock() // Critical section A
			if size < cap {
				size++ // Race: producer and consumer use different mutexes!
				fmt.Printf("Produced! size=%d\n", size)
			}
			mutexA.Unlock()

			// Simulate work
			for i := 0; i < 1000; i++ {
				r = (r + 1) % 10
			}
		}
	}()

	// Consumer section - uses critical section B
	go func() {
		defer wg.Done()
		r := uint(0)

		for {
			mutexB.Lock() // Critical section B - different mutex!
			if size > 0 {
				size-- // Race: producer and consumer use different mutexes!
				fmt.Printf("Consumed! size=%d\n", size)
			}
			mutexB.Unlock()

			// Simulate work
			for i := 0; i < 1000; i++ {
				r = (r + 1) % 10
			}
		}
	}()

	wg.Wait()
}
