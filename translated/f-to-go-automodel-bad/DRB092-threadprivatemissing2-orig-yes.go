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

// Module DRB092 translated to package-level variables
var sum0, sum1 int

func main() {
	var i, sum int
	sum = 0
	sum0 = 0
	sum1 = 0

	//$omp parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()
	chunkSize := 1001 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	//$omp do
	for start := 1; start <= 1001; start += chunkSize {
		end := start + chunkSize - 1
		if end > 1001 {
			end = 1001
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				sum0 = sum0 + i
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

	for i = 1; i <= 1001; i++ {
		sum1 = sum1 + i
	}

	fmt.Printf("sum = %d sum1 = %d\n", sum, sum1)
}
