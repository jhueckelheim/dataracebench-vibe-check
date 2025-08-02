//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Classic i-k-j matrix multiplication. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func main() {
	var N, M, K, length int
	var a, b, c [][]float64

	length = 100
	N = length
	M = length
	K = length

	a = make([][]float64, N)
	b = make([][]float64, M)
	c = make([][]float64, K)
	for i := range a {
		a[i] = make([]float64, M)
	}
	for i := range b {
		b[i] = make([]float64, K)
	}
	for i := range c {
		c[i] = make([]float64, N)
	}

	//$omp parallel do private(j, l)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := N / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= N; start += chunkSize {
		end := start + chunkSize - 1
		if end > N {
			end = N
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				for l := 1; l <= K; l++ {
					for j := 1; j <= M; j++ {
						c[i-1][j-1] = c[i-1][j-1] + a[i-1][l-1]*b[l-1][j-1]
					}
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}
