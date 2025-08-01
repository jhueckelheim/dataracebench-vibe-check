/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
This example is referred from DataRaceOnAccelerator : A Micro-benchmark Suite for Evaluating
Correctness Tools Targeting Accelerators.
Though we have used critical directive to ensure that addition and subtraction are not overlapped,
due to different locks addlock and sublock, operations can interleave each other.
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

	// Two different mutexes - this creates the race condition
	var addLock sync.Mutex
	var subLock sync.Mutex

	// Simulate teams distribute parallel for
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of the iterations
			start := teamID * N / numGoroutines
			end := (teamID + 1) * N / numGoroutines

			for i := start; i < end; i++ {
				// Critical section with addlock
				addLock.Lock()
				variable++ // RACE: Not protected from subLock operations
				addLock.Unlock()

				// Critical section with sublock
				subLock.Lock()
				variable -= 2 // RACE: Not protected from addLock operations
				subLock.Unlock()
			}
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d\n", variable)
}
