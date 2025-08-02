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
	var i, j, n, m int
	var b [][]float64

	n = 100
	m = 100

	b = make([][]float64, n)
	for i := range b {
		b[i] = make([]float64, m)
	}

	for i = 1; i <= n; i++ {
		for j = 1; j <= m; j++ {
			b[i-1][j-1] = float64(i * j)
		}
	}

	for i = 2; i <= n; i++ {
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
			go func(i, start, end int) {
				defer wg.Done()
				for j := start; j <= end; j++ {
					b[i-1][j-1] = b[i-2][j-2]
				}
			}(i, start, end)
		}
		wg.Wait()
		//$omp end parallel do
	}
}
