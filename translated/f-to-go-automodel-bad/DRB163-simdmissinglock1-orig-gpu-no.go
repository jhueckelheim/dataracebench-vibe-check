//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent access of var@26:13 has no atomicity violation. No data race present.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

var var1 [16]int
var i, j int

func main() {
	for i = 1; i <= 16; i++ {
		var1[i-1] = 0
	}

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do reduction(+:var)
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
			localVar := [16]int{}
			for i := start; i <= end; i++ {
				//$omp simd
				for j = 1; j <= 16; j++ {
					localVar[j-1] = localVar[j-1] + 1
				}
				//$omp end simd
			}
			mu.Lock()
			for j = 0; j < 16; j++ {
				var1[j] = var1[j] + localVar[j]
			}
			mu.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	for i = 1; i <= 16; i++ {
		if var1[i-1] != 20 {
			fmt.Printf("%d %d\n", var1[i-1], i)
		}
	}
}
