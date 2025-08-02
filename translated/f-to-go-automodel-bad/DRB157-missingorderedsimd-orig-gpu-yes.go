//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Due to distribute parallel for simd directive at line 23, there is a data race at line 25.
//Data Rae Pairs, var@25:9:W vs. var@25:15:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var var1 [100]int
	var i int

	for i = 1; i <= 100; i++ {
		var1[i-1] = 1
	}

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do simd safelen(16)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 84 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 17; start <= 100; start += chunkSize {
		end := start + chunkSize - 1
		if end > 100 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				var1[i-1] = var1[i-17] + 1
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do simd
	//$omp end target

	fmt.Printf("%d\n", var1[97])
}
