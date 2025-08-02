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

// Module DRB066 translated to package-level functions
func setup(N int) {
	var tar1, tar2 []float64

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
				tar1[i-1] = 0.0
				tar2[i-1] = float64(i) * 2.5
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	// Uncomment to print values for checking
	// fmt.Printf("%f %f\n", tar1[N-1], tar2[N-1])
}

func main() {
	var N int
	N = 1000

	setup(N)
}
