//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The distribute parallel do directive at line 19 will execute loop using multiple teams.
//The loop iterations are distributed across the teams in chunks in round robin fashion.
//The missing lock enclosing var@21 leads to data race. Data Race Pairs, var@21:9:W vs. var@21:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var variable int

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()
	chunkSize := 100 / numTeams

	for team := 0; team < numTeams; team++ {
		start := team*chunkSize + 1
		end := start + chunkSize - 1
		if team == numTeams-1 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				// MISSING: Any synchronization
				variable = variable + 1 // RACE: Multiple teams modify without protection
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", variable)
}