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

// Module DRB085 translated to package-level variables
var sum0, sum1 int64

func foo(i int64) {
	sum0 = sum0 + i
}

func main() {
	var length int
	var i, sum int64
	length = 1000
	sum = 0

	//$omp parallel copyin(sum0)
	var wg sync.WaitGroup
	var mu sync.Mutex
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
				foo(int64(i))
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end do

	//$omp critical
	mu.Lock()
	sum = sum + sum0
	mu.Unlock()
	//$omp end critical
	//$omp end parallel

	for i = 1; i <= length; i++ {
		sum1 = sum1 + i
	}

	fmt.Printf("sum = %d sum1 = %d\n", sum, sum1)
}
