/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Concurrent access on same variable leads to race condition if two different
locks are used. Here we use atomic operations to ensure that addition
and subtraction are not interleaved. No data race pairs.
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const N = 100

func main() {
	var variable int64 = 0
	var wg sync.WaitGroup

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
				// Atomic increment
				atomic.AddInt64(&variable, 1)

				// Atomic decrement by 2
				atomic.AddInt64(&variable, -2)
			}
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d\n", variable)
}
