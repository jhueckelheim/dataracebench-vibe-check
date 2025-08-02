//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//use of omp target + map + array sections derived from pointers. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func foo(a, b []float64, N int) float64 {
	//$omp target map(to:a(1:N)) map(from:b(1:N))
	//$omp parallel do
	// Simulate target parallel execution with goroutines
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
				b[i-1] = a[i-1] * float64(i) // No race - proper partitioning
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
	//$omp end target

	return 0.0 // Function return value
}

func main() {
	var length int
	var a, b []float64
	var x float64

	length = 1000

	a = make([]float64, length)
	b = make([]float64, length)

	// Initialize arrays
	for i := 1; i <= length; i++ {
		a[i-1] = float64(i) / 2.0
		b[i-1] = 0.0
	}

	x = foo(a, b, length)
	fmt.Printf("b(50) = %f\n", b[49])

	// deallocate(a,b)
	_ = x // Use x to avoid unused variable warning
}
