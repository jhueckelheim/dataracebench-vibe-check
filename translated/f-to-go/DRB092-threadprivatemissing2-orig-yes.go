//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A file-scope variable used within a function called by a parallel region.
//No threadprivate is used to avoid data races.
//This is the case for a variable referenced within a construct.
//
//Data race pairs  sum0@34:13:W vs. sum0@34:20:R
//                 sum0@34:13:W vs. sum0@34:13:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variables (module equivalent)
var sum0, sum1 int // sum0 is shared across threads - RACE CONDITION

func main() {
	var sum int
	sum = 0
	sum0 = 0
	sum1 = 0

	//$omp parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp do
			chunkSize := 1001 / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = 1001
			}

			for i := start; i <= end; i++ {
				sum0 = sum0 + i // RACE: Multiple threads modifying shared sum0
			}
			//$omp end do

			//$omp critical
			mu.Lock()
			sum = sum + sum0 // RACE: Reading shared sum0
			mu.Unlock()
			//$omp end critical
		}()
	}
	wg.Wait()
	//$omp end parallel

	// Sequential computation for comparison
	for i := 1; i <= 1001; i++ {
		sum1 = sum1 + i
	}

	fmt.Printf("sum = %d sum1 = %d\n", sum, sum1)
}
