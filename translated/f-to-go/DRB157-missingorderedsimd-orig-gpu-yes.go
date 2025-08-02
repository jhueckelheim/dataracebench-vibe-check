//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Due to distribute parallel for simd directive at line 23, there is a data race at line 25.
//Data Race Pairs, var@25:9:W vs. var@25:15:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var variable [100]int

	// Initialize
	for i := 1; i <= 100; i++ {
		variable[i-1] = 1
	}

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do simd safelen(16)
	// SIMD with safelen(16) but distance is 16 - creates race
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()
	chunkSize := 84 / numTeams // 17 to 100 = 84 iterations

	for team := 0; team < numTeams; team++ {
		start := team*chunkSize + 17
		end := start + chunkSize - 1
		if team == numTeams-1 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				variable[i-1] = variable[i-17] + 1 // RACE: Distance 16 = safelen, creates race
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do simd
	//$omp end target

	fmt.Printf("%d\n", variable[97])
}