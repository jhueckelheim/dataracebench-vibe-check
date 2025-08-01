/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Concurrent access of var in an intra region. Lock ensures that there is no data race.
Uses num_teams(1) to ensure all work is in one team, making the lock effective.
*/

package main

import (
	"fmt"
	"sync"
)

const N = 100

func main() {
	var variable int = 0
	var lock sync.Mutex
	var wg sync.WaitGroup

	// Simulate single team (num_teams(1)) with parallel execution
	// All goroutines are in the same "team" so the lock protects properly
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Each thread processes a chunk of the iterations
			start := threadID * N / numGoroutines
			end := (threadID + 1) * N / numGoroutines

			for i := start; i < end; i++ {
				// Lock protects all threads in the single team
				lock.Lock()
				variable++ // SAFE: Protected by single lock across all threads
				lock.Unlock()
			}
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d\n", variable)
}
