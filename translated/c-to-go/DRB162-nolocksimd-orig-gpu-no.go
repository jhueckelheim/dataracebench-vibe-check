/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Reduction clause will ensure there is no data race in var. No Data Race.
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	N = 20
	C = 8
)

func main() {
	var variable [C]int64
	var wg sync.WaitGroup

	// Initialize array
	for i := 0; i < C; i++ {
		variable[i] = 0
	}

	// Simulate teams distribute parallel for with reduction
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of N iterations
			start := teamID * N / numGoroutines
			end := (teamID + 1) * N / numGoroutines

			// Local reduction array
			localVar := make([]int64, C)

			for i := start; i < end; i++ {
				// SIMD loop - local accumulation (reduction pattern)
				for j := 0; j < C; j++ {
					localVar[j]++
				}
			}

			// Atomic reduction to global array
			for j := 0; j < C; j++ {
				atomic.AddInt64(&variable[j], localVar[j])
			}
		}(t)
	}

	wg.Wait()

	// Check results
	for i := 0; i < C; i++ {
		if variable[i] != N {
			fmt.Printf("%d \n", variable[i])
		}
	}
}
