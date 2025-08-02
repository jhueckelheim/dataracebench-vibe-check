//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Only the outmost loop can be parallelized. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func foo() {
	var n, m, length int
	var b [][]float64

	length = 100
	b = make([][]float64, length)
	for i := range b {
		b[i] = make([]float64, length)
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
				for j := 1; j <= m-1; j++ {
					b[i-1][j-1] = b[i-1][j]
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
