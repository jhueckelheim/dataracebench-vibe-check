/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
The var is atomic update. Hence, there is no data race pair.
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

	// Simulate teams distribute with atomic updates
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of the iterations
			start := teamID * N / numGoroutines
			end := (teamID + 1) * N / numGoroutines

			for i := start; i < end; i++ {
				// Atomic update equivalent to #pragma omp atomic update
				atomic.AddInt64(&variable, 1)
			}
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d \n", variable)
}
