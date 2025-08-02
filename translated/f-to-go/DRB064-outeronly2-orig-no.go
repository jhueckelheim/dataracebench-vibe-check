//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Only the outmost loop can be parallelized.
//The inner loop has loop carried true data dependence.
//However, the loop is not parallelized so no race condition.

package main

import (
	"runtime"
	"sync"
)

func foo() {
	var n, m, length int
	var b [][]float32

	length = 100
	b = make([][]float32, length)
	for i := 0; i < length; i++ {
		b[i] = make([]float32, length)
	}
	n = length
	m = length

	//$omp parallel do private(j)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := n / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= n; start += chunkSize {
		end := start + chunkSize - 1
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				for j := 2; j <= m; j++ { // j is private, inner loop has dependence but is not parallelized
					b[i-1][j-1] = b[i-1][j-2] // No race - proper data partitioning
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}

func main() {
	foo()
}
