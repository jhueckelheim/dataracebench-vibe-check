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
	"runtime"
	"sync"
)

func main() {
	var var1 [8]int
	var i, j int

	for i = 1; i <= 8; i++ {
		var1[i-1] = 0
	}

	//$omp target map(tofrom:var) device(0)
	//$omp teams num_teams(1) thread_limit(1048)
	//$omp distribute parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 20 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= 20; start += chunkSize {
		end := start + chunkSize - 1
		if end > 20 {
			end = 20
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				//$omp simd
				for j = 1; j <= 8; j++ {
					var1[j-1] = var1[j-1] + 1
				}
				//$omp end simd
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end distribute parallel do
	//$omp end teams
	//$omp end target

	fmt.Printf("%d\n", var1[7])
}
