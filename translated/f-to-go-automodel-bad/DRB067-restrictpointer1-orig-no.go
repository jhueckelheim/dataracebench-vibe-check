//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Array initialization using assignments. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Module DRB067 translated to package-level functions
func foo(newSxx, newSyy []float64, length int) {
	var tar1, tar2 []float64

	newSxx = make([]float64, length)
	newSyy = make([]float64, length)
	tar1 = make([]float64, length)
	tar2 = make([]float64, length)

	newSxx = tar1[:]
	newSyy = tar2[:]

	//$omp parallel do private (i) firstprivate (len)
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
		go func(start, end, length int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				tar1[i-1] = 0.0
				tar2[i-1] = 0.0
			}
		}(start, end, length)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("%f %f\n", tar1[length-1], tar2[length-1])
}

func main() {
	var length int = 1000
	var newSxx, newSyy []float64

	newSxx = make([]float64, length)
	newSyy = make([]float64, length)

	foo(newSxx, newSyy, length)
}
