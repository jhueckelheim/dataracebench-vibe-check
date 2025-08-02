//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Classic PI calculation using reduction. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var interval_width, pi float64
	var num_steps int64

	pi = 0.0
	num_steps = 2000000000
	interval_width = 1.0 / float64(num_steps)

	//$omp parallel do reduction(+:pi) private(x)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()
	chunkSize := int(num_steps) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := int64(1); start <= num_steps; start += int64(chunkSize) {
		end := start + int64(chunkSize) - 1
		if end > num_steps {
			end = num_steps
		}
		wg.Add(1)
		go func(start, end int64) {
			defer wg.Done()
			localPi := 0.0
			for i := start; i <= end; i++ {
				x := (float64(i) + 0.5) * interval_width // x is private to each goroutine
				localPi = localPi + 1.0/(x*x+1.0)
			}
			// Reduction
			mu.Lock()
			pi += localPi
			mu.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	pi = pi * 4.0 * interval_width
	fmt.Printf("PI = %24.20f\n", pi)
}
