//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Matrix-vector multiplication: outer-level loop parallelization. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func foo() {
	var N int
	var sum float64
	var a [][]float64
	var v, v_out []float64

	N = 100
	a = make([][]float64, N)
	for i := range a {
		a[i] = make([]float64, N)
	}
	v = make([]float64, N)
	v_out = make([]float64, N)

	//$omp parallel do private (i,j,sum)
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
				sum = 0.0
				for j := 1; j <= N; j++ {
					sum = sum + a[i-1][j-1]*v[j-1]
				}
				v_out[i-1] = sum
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}

func main() {
	foo()
}
