/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* This kernel imitates the nature of a program from the NAS Parallel Benchmarks 3.0 MG suit.
 * Use of private clause will ensure that there is no data race. No Data Race Pairs.
 */

package main

import (
	"sync"
)

func main() {
	var a [12][12][12]float64
	var wg sync.WaitGroup

	m := 3.0

	// Simulate parallel for private(j,k,tmp1)
	numGoroutines := 4
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Each thread processes a chunk of i values
			start := threadID * 12 / numGoroutines
			end := (threadID + 1) * 12 / numGoroutines

			for i := start; i < end; i++ {
				// Private variables j, k, tmp1 (local to each goroutine)
				for j := 0; j < 12; j++ {
					for k := 0; k < 12; k++ {
						tmp1 := 6.0 / m
						a[i][j][k] = tmp1 + 4
					}
				}
			}
		}(t)
	}

	wg.Wait()
}
