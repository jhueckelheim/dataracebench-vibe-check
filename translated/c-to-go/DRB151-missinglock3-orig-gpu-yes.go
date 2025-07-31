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
The distribute parallel for directive will execute loop using multiple teams.
The loop iterations are distributed across the teams in chunks in round robin fashion.
The missing lock enclosing var leads to data race.
Data Race Pairs, var:W vs. var:W
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var variable int = 0
	var wg sync.WaitGroup

	// Simulate teams distribute parallel for with no protection
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of the 100 iterations
			start := teamID * 100 / numGoroutines
			end := (teamID + 1) * 100 / numGoroutines

			for i := start; i < end; i++ {
				// NO PROTECTION - Direct data race
				variable++ // RACE: Concurrent access without synchronization
			}
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d\n", variable)
}
