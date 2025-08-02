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
	var i, error, length, b int
	var a []int

	length = 1000
	b = 5
	a = make([]int, length)

	for i = 1; i <= length; i++ {
		a[i-1] = i
	}

	//$omp parallel shared(b, error)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := length / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	//$omp do
	for start := 1; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = b + a[i-1]*5
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end do nowait
	//$omp end parallel

	//$omp barrier
	//$omp single
	error = a[8] + 1
	//$omp end single

	fmt.Printf("error =%8d\n", error)
}
