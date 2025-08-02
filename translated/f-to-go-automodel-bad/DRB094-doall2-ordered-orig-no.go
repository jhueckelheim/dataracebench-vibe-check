//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimensional array computation:
//ordered(2) is used to associate two loops with omp for.
//The corresponding loop iteration variables are private.
//
//ordered(n) is an OpenMP 4.5 addition. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Module DRB094 translated to package-level variables
var a [][]int

func main() {
	var length int
	length = 100

	a = make([][]int, length)
	for i := range a {
		a[i] = make([]int, length)
	}

	//$omp parallel do ordered(2)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	totalIterations := length * length
	chunkSize := totalIterations / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 0; start < totalIterations; start += chunkSize {
		end := start + chunkSize - 1
		if end >= totalIterations {
			end = totalIterations - 1
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for iter := start; iter <= end; iter++ {
				i := iter/length + 1
				j := iter%length + 1
				a[i-1][j-1] = a[i-1][j-1] + 1
				//$omp ordered depend(sink:i-1,j) depend (sink:i,j-1)
				fmt.Printf("test i = %d   j = %d\n", i, j)
				//$omp ordered depend(source)
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}
