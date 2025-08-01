/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

no race
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
	mutex    sync.Mutex
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
					mutex.Lock() // Shared critical section
					if size < cap {
						size++          // produce
						localPackages-- // produced a package
						fmt.Printf("Producer %d produced! size=%d\n", id, size)
					}
					mutex.Unlock()
				}
			} else {
				// I am a consumer
				for localPackages > 0 {
					mutex.Lock() // Same shared critical section as producer
					if size > 0 {
						size--          // consume
						localPackages-- // consumed a package
						fmt.Printf("Consumer %d consumed! size=%d\n", id-nprod, size)
					}
					mutex.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
}
