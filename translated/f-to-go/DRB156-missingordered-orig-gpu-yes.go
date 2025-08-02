//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Missing ordered directive causes data race pairs var@24:9:W vs. var@24:18:R

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
	//$omp teams distribute parallel do
	// MISSING: ordered directive - causes race
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()
	chunkSize := 99 / numTeams // 2 to 100 = 99 iterations

	for team := 0; team < numTeams; team++ {
		start := team*chunkSize + 2
		end := start + chunkSize - 1
		if team == numTeams-1 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				variable[i-1] = variable[i-2] + 1 // RACE: Reading/writing without ordering
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", variable[99])
}