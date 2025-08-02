//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Arrays passed as function parameters. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

// Module DRB050 translated to package-level variables
var o1, c []float64

func foo1(o1, c []float64, length int) {
	//$omp parallel do private(volnew_o8)
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
			for i := start; i <= end; i++ {
				volnew_o8 := 0.5 * c[i-1]
				o1[i-1] = volnew_o8
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
	// Uncomment to print a value for checking
	// fmt.Printf("%f\n", o1[49])
}

func main() {
	o1 = make([]float64, 100)
	c = make([]float64, 100)

	foo1(o1, c, 100)
}
