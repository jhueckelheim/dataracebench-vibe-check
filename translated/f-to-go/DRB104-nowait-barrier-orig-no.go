//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This example is based on one code snippet extracted from a paper:
//Ma etc. Symbolic Analysis of Concurrency Errors in OpenMP Programs, ICPP 2013
//
//Explicit barrier to counteract nowait. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var err, length, b int
	var a []int

	length = 1000
	b = 5
	a = make([]int, length)

	// Initialize array
	for i := 1; i <= length; i++ {
		a[i-1] = i
	}

	//$omp parallel shared(b, error)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp do
			chunkSize := length / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = length
			}

			for i := start; i <= end; i++ {
				a[i-1] = b + a[i-1]*5 // No race - proper partitioning
			}
			//$omp end do nowait
			// nowait: don't wait here, but explicit barrier below ensures synchronization
		}()
	}
	wg.Wait() // This acts as the explicit barrier
	//$omp end parallel

	//$omp barrier (already handled by wg.Wait())
	//$omp single
	// Only executed once after barrier
	err = a[8] + 1 // a(9) in Fortran is a[8] in Go
	//$omp end single

	fmt.Printf("error = %8d\n", err)

	// deallocate(a)
}
