/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
This example is referred from DRACC by Adrian Schmitz et al.
Missing ordered directive causes data race pairs var[i]:W vs. var[i-1]:R
*/

package main

import (
	"fmt"
	"sync"
)

const N = 100

func main() {
	var variable [N]int

	// Initialize array
	for i := 0; i < N; i++ {
		variable[i] = 0
	}

	// Simulate teams distribute parallel for WITHOUT ordering
	var wg sync.WaitGroup
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of iterations
			start := teamID*(N-1)/numGoroutines + 1
			end := (teamID+1)*(N-1)/numGoroutines + 1
			if end > N {
				end = N
			}

			for i := start; i < end; i++ {
				// RACE: Reading var[i-1] while another goroutine might be writing it
				variable[i] = variable[i-1] + 1
			}
		}(t)
	}

	wg.Wait()

	// Check results
	for i := 0; i < N; i++ {
		if variable[i] != i {
			fmt.Printf("Data Race Present\n")
			return
		}
	}

	fmt.Printf("No race detected in this run\n")
}
