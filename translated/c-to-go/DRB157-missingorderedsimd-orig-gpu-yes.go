/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
This kernel is modified version from "DataRaceOnAccelerator A Micro-benchmark Suite for Evaluating
Correctness Tools Targeting Accelerators" by Adrian Schmitz et al.
Due to distribute parallel for simd directive, there is a data race.
Data Race Pairs, var[i]:W vs. var[i-C]:R
*/

package main

import (
	"fmt"
	"sync"
)

const (
	N = 100
	C = 16
)

func main() {
	var variable [N]int

	// Initialize array
	for i := 0; i < N; i++ {
		variable[i] = 0
	}

	// Simulate teams distribute parallel for simd with safelen(C)
	// SIMD with safelen(C) means dependencies separated by < C are violated
	var wg sync.WaitGroup
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of iterations from C to N
			start := teamID*(N-C)/numGoroutines + C
			end := (teamID+1)*(N-C)/numGoroutines + C
			if end > N {
				end = N
			}

			for i := start; i < end; i++ {
				// RACE: Dependencies within safelen(C=16) are violated
				// Reading var[i-C] while potentially writing to nearby indices
				variable[i] = variable[i-C] + 1
			}
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d\n", variable[97])
}
