//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A file-scope variable used within a function called by a parallel region.
//No threadprivate is used to avoid data races.
//
//Data race pairs  sum@39:13:W vs. sum@39:19:R
//                 sum@39:13:W vs. sum@39:13:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variables (module equivalent)
var sum0, sum1 int64 // sum0 is shared across threads - RACE CONDITION

func foo(i int64) {
	sum0 = sum0 + i // RACE: Multiple threads accessing shared sum0
}

func main() {
	var sum int64
	sum = 0

	//$omp parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp do
			// Simulate work distribution
			chunkSize := 1001 / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = 1001 // Last thread handles remainder
			}

			for i := int64(start); i <= int64(end); i++ {
				foo(i) // RACE: All threads call foo which modifies shared sum0
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
	for i := int64(1); i <= 1001; i++ {
		sum1 = sum1 + i
	}

	fmt.Printf("sum = %d sum1 = %d\n", sum, sum1)
}
