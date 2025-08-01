/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
This kernel is referred from "DataRaceOnAccelerator A Micro-benchmark Suite for Evaluating
Correctness Tools Targeting Accelerators" by Adrian Schmitz et al.
Concurrent access of var in an intra region. Missing Lock leads to intra region data race.
Data Race pairs, var:W vs. var:W
*/

package main

import (
	"fmt"
	"sync"
)

const N = 100

func main() {
	var variable int = 0
	var wg sync.WaitGroup

	// Simulate single team (num_teams(1)) but missing lock protection
	// Even in single team, parallel threads cause race condition
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Each thread processes a chunk of the iterations
			start := threadID * N / numGoroutines
			end := (threadID + 1) * N / numGoroutines

			for i := start; i < end; i++ {
				// NO LOCK PROTECTION - Race condition
				variable++ // RACE: Concurrent access without synchronization
			}
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d \n", variable)
}
