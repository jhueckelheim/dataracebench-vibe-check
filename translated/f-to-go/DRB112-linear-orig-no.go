//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//omp for loop is allowed to use the linear clause, an OpenMP 4.5 addition. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func main() {
	var length, i int
	var a, b, c []float64

	length = 100
	i = 0

	a = make([]float64, length)
	b = make([]float64, length)
	c = make([]float64, length)

	// Initialize arrays
	for i = 1; i <= length; i++ {
		a[i-1] = float64(i) / 2.0
		b[i-1] = float64(i) / 3.0
		c[i-1] = float64(i) / 7.0
	}

	//$omp parallel do linear(j)
	// With linear(j): each thread gets a private j that increments predictably
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := length / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			// linear(j): j starts at its initial value for each thread
			j := start - 1 // Each thread gets proper j value
			for i := start; i <= end; i++ {
				c[j] = c[j] + a[i-1]*b[i-1] // No race - j is private and linear
				j = j + 1                   // No race - j increments linearly per thread
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	// print*,'c(50) =',c(50)

	// if(allocated(a))deallocate(a)
	// if(allocated(b))deallocate(b)
}
