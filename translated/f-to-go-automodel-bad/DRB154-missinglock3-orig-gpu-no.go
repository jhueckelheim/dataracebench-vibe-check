//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent accessing var@25:9 may cause atomicity violation and inter region data race.
//Lock and reduction clause at line 22, avoids this. No Data Race Pair.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var lck sync.Mutex
	var var1, i int
	var1 = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute reduction(+:var)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()
	chunkSize := 100 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= 100; start += chunkSize {
		end := start + chunkSize - 1
		if end > 100 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localVar := 0
			for i := start; i <= end; i++ {
				lck.Lock()
				localVar = localVar + 1
				lck.Unlock()
			}
			mu.Lock()
			var1 = var1 + localVar
			mu.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute
	//$omp end target

	fmt.Printf("%d\n", var1)
}
