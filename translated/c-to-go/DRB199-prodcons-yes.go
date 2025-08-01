/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

race introduced because critical sections have different names for producer and consumer.
Data race pair: size@34:11:W vs. size@45:11:W
*/

package main

import (
	"fmt"
	"sync"
)

var (
	nprod    = 4
	ncons    = 4
	cap      = 5
	size     = 0
	packages = 1000
	mutexA   sync.Mutex // Critical section A
	mutexB   sync.Mutex // Critical section B
)

func main() {
	nthread := nprod + ncons
	var wg sync.WaitGroup

	for i := 0; i < nthread; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			localPackages := packages

			if id < nprod {
				// I am a producer
				for localPackages > 0 {
					mutexA.Lock() // Critical section A
					if size < cap {
						size++          // Race: different mutex than consumer!
						localPackages-- // produced a package
						fmt.Printf("Producer %d produced! size=%d\n", id, size)
					}
					mutexA.Unlock()
				}
			} else {
				// I am a consumer
				for localPackages > 0 {
					mutexB.Lock() // Critical section B - different mutex!
					if size > 0 {
						size--          // Race: different mutex than producer!
						localPackages-- // consumed a package
						fmt.Printf("Consumer %d consumed! size=%d\n", id-nprod, size)
					}
					mutexB.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
}
