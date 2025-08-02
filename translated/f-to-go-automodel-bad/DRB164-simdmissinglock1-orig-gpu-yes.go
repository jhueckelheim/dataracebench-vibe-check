//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent access of var@30:13 without acquiring locks causes atomicity violation. Data race present.
//Data Race Pairs, var@30:13:W vs. var@30:22:R

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
	//$omp teams distribute parallel do
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
				for j = 1; j <= 16; j++ {
					var1[j-1] = var1[j-1] + 1
				}
				//$omp end simd
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", var1[15])
}
