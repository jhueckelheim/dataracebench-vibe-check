//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimensional array computation using loops: missing private(j).
//References to j in the loop cause data races.
//Data race pairs (we allow multiple ones to preserve the pattern):
//  Write_set = {j@28:12} (implicit step by +1)
//  Read_set = {j@29:17, j@29:26, j@28:12} (implicit step by +1)
//  Any pair from Write_set vs. Write_set  and Write_set vs. Read_set is a data race pair.

package main

import (
	"runtime"
	"sync"
)

func main() {
	var i, j, length int // j is shared among all threads - RACE CONDITION
	var a [][]int

	length = 100
	a = make([][]int, length)
	for i = 0; i < length; i++ {
		a[i] = make([]int, length)
	}

	//$omp parallel do
	var wg sync.WaitGroup
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
			for i := start; i <= end; i++ {
				for j = 1; j <= 100; j++ { // RACE: j is shared, not private!
					a[i-1][j-1] = a[i-1][j-1] + 1
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	// deallocate(a)
}
