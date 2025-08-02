//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Freshly allocated pointers do not alias to each other. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func setup(N int) {
	var tar1, tar2 []float64

	// Allocate arrays (freshly allocated, no aliasing)
	tar1 = make([]float64, N)
	tar2 = make([]float64, N)

	//$omp parallel do schedule(static)
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
				tar1[i-1] = 0.0              // No aliasing - tar1 and tar2 are separate arrays
				tar2[i-1] = float64(i) * 2.5 // No race condition
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	// print*,tar1(N),tar2(N)
	// In Go, slices are automatically garbage collected when out of scope
}

func main() {
	N := 1000
	setup(N)
}
