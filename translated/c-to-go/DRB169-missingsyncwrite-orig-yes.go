/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* This kernel imitates the nature of a program from the NAS Parallel Benchmarks 3.0 MG suit.
 * Due to missing construct to write r1[k] synchronously, there is a Data Race.
 * Data Race Pair, r1[k]:W vs. r1[k]:W
 * */

package main

import (
	"fmt"
	"sync"
)

const N = 8

func main() {
	var r1 [N]float64
	var r [N][N][N]float64
	var wg sync.WaitGroup

	// Initialize 3D array
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			for k := 0; k < N; k++ {
				r[i][j][k] = float64(i)
			}
		}
	}

	// Simulate parallel for default(shared) private(j,k)
	numGoroutines := 4
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Each thread processes a chunk of i values
			start := threadID*(N-2)/numGoroutines + 1
			end := (threadID+1)*(N-2)/numGoroutines + 1
			if end > N-1 {
				end = N - 1
			}

			// Private variables j, k (automatic in Go closures)
			for i := start; i < end; i++ {
				for j := 1; j < N-1; j++ {
					for k := 0; k < N; k++ {
						// RACE: Multiple goroutines write to same r1[k] without synchronization
						r1[k] = r[i][j-1][k] + r[i][j+1][k] + r[i-1][j][k] + r[i+1][j][k]
					}
				}
			}
		}(t)
	}

	wg.Wait()

	// Print results
	for k := 0; k < N; k++ {
		fmt.Printf("%f ", r1[k])
	}
	fmt.Printf("\n")
}
