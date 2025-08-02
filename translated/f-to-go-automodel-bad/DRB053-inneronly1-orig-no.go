//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Example with loop-carried data dependence at the outer level loop.
//But the inner level loop can be parallelized.

package main

import (
	"runtime"
	"sync"
)

func main() {
	var a [][]float64

	a = make([][]float64, 20)
	for i := range a {
		a[i] = make([]float64, 20)
	}

	for i := 1; i <= 20; i++ {
		for j := 1; j <= 20; j++ {
			a[i-1][j-1] = 0.0
		}
	}

	for i := 1; i <= 19; i++ {
		//$omp parallel do
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := 20 / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 1; start <= 20; start += chunkSize {
			end := start + chunkSize - 1
			if end > 20 {
				end = 20
			}
			wg.Add(1)
			go func(i, start, end int) {
				defer wg.Done()
				for j := start; j <= end; j++ {
					a[i-1][j-1] = a[i-1][j-1] + a[i][j-1]
				}
			}(i, start, end)
		}
		wg.Wait()
		//$omp end parallel do
	}
}
