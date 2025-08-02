//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Reduction clause at line 23:34 will ensure there is no data race in var@27:13. No Data Race.

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
	//$omp distribute parallel do reduction(+:var)
	var wg sync.WaitGroup
	var mu sync.Mutex
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
			localVar := [8]int{}
			for i := start; i <= end; i++ {
				//$omp simd
				for j = 1; j <= 8; j++ {
					localVar[j-1] = localVar[j-1] + 1
				}
				//$omp end simd
			}
			mu.Lock()
			for j = 0; j < 8; j++ {
				var1[j] = var1[j] + localVar[j]
			}
			mu.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end distribute parallel do
	//$omp end teams
	//$omp end target

	for i = 1; i <= 8; i++ {
		if var1[i-1] != 20 {
			fmt.Printf("%d\n", var1[i-1])
		}
	}
}
