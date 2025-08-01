/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

single producer single consumer with critical sections
*/

package main

import (
	"fmt"
	"sync"
)

var (
	cap      = 10
	size     = 0
	packages = 1000
	mutex    sync.Mutex
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Producer section
	go func() {
		defer wg.Done()
		localPackages := packages
		r := uint(0)

		for localPackages > 0 {
			mutex.Lock()
			if size < cap {
				size++          // produce
				localPackages-- // produced a package
				fmt.Printf("Produced! size=%d\n", size)
			}
			mutex.Unlock()

			// Simulate work
			for i := 0; i < 1000; i++ {
				r = (r + 1) % 10
			}
		}
	}()

	// Consumer section
	go func() {
		defer wg.Done()
		localPackages := packages
		r := uint(0)

		for localPackages > 0 {
			mutex.Lock()
			if size > 0 {
				size--          // consume
				localPackages-- // consumed a package
				fmt.Printf("Consumed! size=%d\n", size)
			}
			mutex.Unlock()

			// Simulate work
			for i := 0; i < 1500; i++ {
				r = (r + 1) % 10
			}
		}
	}()

	wg.Wait()
}
