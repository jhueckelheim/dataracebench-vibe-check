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
Concurrent access of var without acquiring locks causes atomicity violation. Data race present.
Data Race Pairs, var[i]:W vs. var[i]:W
*/

package main

import (
	"fmt"
	"sync"
)

const (
	N = 100
	C = 64
)

func main() {
	var variable [C]int
	var wg sync.WaitGroup

	// Initialize array
	for i := 0; i < C; i++ {
		variable[i] = 0
	}

	// Simulate teams distribute parallel for WITHOUT reduction
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of N iterations
			start := teamID * N / numGoroutines
			end := (teamID + 1) * N / numGoroutines

			for i := start; i < end; i++ {
				// SIMD loop - direct updates without protection
				for j := 0; j < C; j++ {
					// RACE: Multiple goroutines increment same array elements
					variable[j]++ // RACE: Concurrent access without synchronization
				}
			}
		}(t)
	}

	wg.Wait()

	fmt.Printf("%d\n", variable[63])
}
