//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Matrix-vector multiplication: inner level parallelization. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func foo() {
	var N int
	var a [][]float32
	var v, v_out []float32

	N = 1000
	a = make([][]float32, N)
	for i := 0; i < N; i++ {
		a[i] = make([]float32, N)
	}
	v = make([]float32, N)
	v_out = make([]float32, N)

	for i := 1; i <= N; i++ {
		sum := float32(0.0)

		//$omp parallel do reduction(+:sum)
		var wg sync.WaitGroup
		var mu sync.Mutex
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
				localSum := float32(0.0)
				for j := start; j <= end; j++ {
					localSum = localSum + a[i-1][j-1]*v[j-1]
					fmt.Printf("%f\n", localSum) // Note: this will be printed multiple times from different goroutines
				}
				// Reduction
				mu.Lock()
				sum += localSum
				mu.Unlock()
			}(start, end)
		}
		wg.Wait()
		//$omp end parallel do

		v_out[i-1] = sum
	}
}

func main() {
	foo()
}
