//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//By utilizing the ordered construct @23 the execution will be sequentially consistent.
//No Data Race Pair.

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
	//$omp parallel do ordered
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()
	chunkSize := 99 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 2; start <= 100; start += chunkSize {
		end := start + chunkSize - 1
		if end > 100 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				//$omp ordered
				mu.Lock()
				var1[i-1] = var1[i-2] + 1
				mu.Unlock()
				//$omp end ordered
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
	//$omp end target

	for i = 1; i <= 100; i++ {
		if var1[i-1] != i {
			fmt.Printf("Data Race Present\n")
		}
	}
}
