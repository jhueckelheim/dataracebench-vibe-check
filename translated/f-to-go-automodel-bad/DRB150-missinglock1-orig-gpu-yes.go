//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The distribute parallel do directive at line 22 will execute loop using multiple teams.
//The loop iterations are distributed across the teams in chunks in round robin fashion.
//The omp lock is only guaranteed for a contention group, i.e, within a team.
//Data Race Pair, var@25:9:W vs. var@25:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var var1, i int
	var lck sync.Mutex

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 10 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= 10; start += chunkSize {
		end := start + chunkSize - 1
		if end > 10 {
			end = 10
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				lck.Lock()
				var1 = var1 + 1
				lck.Unlock()
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", var1)
}
