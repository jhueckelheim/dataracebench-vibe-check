/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
The increment operation is team specific as each team works on their individual var.
No Data Race Pair - uses atomic reduction for safe parallel increments.
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

	// Simulate teams distribute parallel for with reduction
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of the work
			start := teamID * N / numGoroutines
			end := (teamID + 1) * N / numGoroutines

			// Local accumulation to reduce contention
			localSum := int64(0)
			for i := start; i < end; i++ {
				localSum++
			}

			// Atomic addition of local sum (reduction pattern)
			atomic.AddInt64(&variable, localSum)
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d\n", variable)
}
