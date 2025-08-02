//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This example is from DRACC by Adrian Schmitz et al.
//Concurrent access on a counter with no lock with simd. Atomicity Violation. Intra Region.
//Data Race Pairs: var@29:13:W vs. var@29:13:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var variable [8]int

	// Initialize
	for i := 1; i <= 8; i++ {
		variable[i-1] = 0
	}

	//$omp target map(tofrom:var) device(0)
	//$omp teams num_teams(1) thread_limit(1048)
	//$omp distribute parallel do
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			//$omp simd
			// SIMD without proper synchronization
			for j := 1; j <= 8; j++ {
				variable[j-1] = variable[j-1] + 1 // RACE: Multiple SIMD lanes accessing same memory
			}
			//$omp end simd
		}()
	}
	wg.Wait()
	//$omp end distribute parallel do
	//$omp end teams
	//$omp end target

	fmt.Printf("%d\n", variable[7])
}