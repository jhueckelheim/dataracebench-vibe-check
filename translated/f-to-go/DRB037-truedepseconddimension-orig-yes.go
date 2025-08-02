//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Only the outmost loop can be parallelized in this program.
//The inner loop has true dependence.
//Data race pair: b[i][j]@29:13:W vs. b[i][j-1]@29:22:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, n, m, length int
	var b [][]float32

	length = 1000
	n = length
	m = length

	b = make([][]float32, length)
	for i = 0; i < length; i++ {
		b[i] = make([]float32, length)
	}

	for i = 1; i <= n; i++ {
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
					b[i-1][j-1] = b[i-1][j-2] // True dependence on second dimension
				}
			}(start, end)
		}
		wg.Wait()
		//$omp end parallel do
	}

	fmt.Printf("b(500,500) = %20.6f\n", b[499][499])

	// deallocate(b)
}