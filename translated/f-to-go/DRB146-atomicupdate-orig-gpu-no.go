//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The var@22:13 is atomic update. Hence, there is no data race.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	var variable int64 // Use int64 for atomic operations
	variable = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute
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
				//$omp atomic update
				atomic.AddInt64(&variable, 1) // No race - atomic operation
				//$omp end atomic
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute
	//$omp end target

	fmt.Printf("%d\n", variable)
}