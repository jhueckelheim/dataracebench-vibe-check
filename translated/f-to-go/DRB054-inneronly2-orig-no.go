//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Example with loop-carried data dependence at the outer level loop.
//The inner level loop can be parallelized. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func main() {
	var n, m int
	var b [][]float32

	n = 100
	m = 100

	b = make([][]float32, n)
	for i := 0; i < n; i++ {
		b[i] = make([]float32, m)
	}

	// Initialize array
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			b[i-1][j-1] = float32(i * j)
		}
	}

	// Outer loop has data dependence (b[i] depends on b[i-1])
	// but inner loop can be parallelized
	for i := 2; i <= n; i++ {
		//$omp parallel do
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := (m - 1) / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 2; start <= m; start += chunkSize {
			end := start + chunkSize - 1
			if end > m {
				end = m
			}
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for j := start; j <= end; j++ {
					b[i-1][j-1] = b[i-2][j-2] // No race on inner loop
				}
			}(start, end)
		}
		wg.Wait()
		//$omp end parallel do
	}

	// deallocate(b)
}
