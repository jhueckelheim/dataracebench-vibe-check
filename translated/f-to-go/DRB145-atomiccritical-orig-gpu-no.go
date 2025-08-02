//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The increment operation at line@22:17 is team specific as each team work on their individual var.
//No Data Race Pair

package main

import (
	"runtime"
	"sync"
)

func main() {
	var variable int
	variable = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do reduction(+:var)
	// Proper reduction across teams
	var wg sync.WaitGroup
	var mu sync.Mutex
	numTeams := runtime.NumCPU()
	chunkSize := 200 / numTeams

	for team := 0; team < numTeams; team++ {
		start := team*chunkSize + 1
		end := start + chunkSize - 1
		if team == numTeams-1 {
			end = 200
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localVar := 0
			for i := start; i <= end; i++ {
				if localVar < 101 {
					localVar = localVar + 1 // No race - each team has private var
				}
			}
			// Reduction
			mu.Lock()
			variable += localVar
			mu.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target
}