//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The increment at line number 22 is critical for the variable
//var@22:13. Therefore, there is a possible Data Race pair var@22:13:W vs. var@22:19:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var var1, i int
	var1 = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 200 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= 200; start += chunkSize {
		end := start + chunkSize - 1
		if end > 200 {
			end = 200
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				//$omp critical
				var1 = var1 + 1
				//$omp end critical
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", var1)
}
