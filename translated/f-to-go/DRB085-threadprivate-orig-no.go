//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A file-scope variable used within a function called by a parallel region.
//Use threadprivate to avoid data races. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable
var sum1 int64

func foo(i int64, localSum0 *int64) {
	*localSum0 = *localSum0 + i // No race - each thread has its own copy
}

func main() {
	var length int
	var sum int64

	length = 1000
	sum = 0

	//$omp parallel copyin(sum0)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Each goroutine has its own sum0 (simulating threadprivate)
			var sum0 int64 = 0 // threadprivate equivalent - each thread has its own copy

			//$omp do
			// Simulate work distribution
			chunkSize := length / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = length // Last thread handles remainder
			}

			for i := int64(start); i <= int64(end); i++ {
				foo(i, &sum0) // No race - each thread has private sum0
			}
			//$omp end do

			//$omp critical
			mu.Lock()
			sum = sum + sum0 // No race - reading private sum0
			mu.Unlock()
			//$omp end critical
		}()
	}
	wg.Wait()
	//$omp end parallel

	// Sequential computation for comparison
	for i := int64(1); i <= int64(length); i++ {
		sum1 = sum1 + i
	}

	fmt.Printf("sum = %d sum1 = %d\n", sum, sum1)
}
