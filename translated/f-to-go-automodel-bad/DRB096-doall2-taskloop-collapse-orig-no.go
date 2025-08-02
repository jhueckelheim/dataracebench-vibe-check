//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimensional array computation:
//Two loops are associated with omp taskloop due to collapse(2).
//Both loop index variables are private.
//taskloop requires OpenMP 4.5 compilers. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Module DRB096 translated to package-level variables
var a [][]int

func main() {
	var length int
	length = 100

	a = make([][]int, length)
	for i := range a {
		a[i] = make([]int, length)
	}

	//$omp parallel
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp single
		var wgSingle sync.WaitGroup
		numCPU := runtime.NumCPU()
		totalIterations := length * length
		chunkSize := totalIterations / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		//$omp taskloop collapse(2)
		for start := 0; start < totalIterations; start += chunkSize {
			end := start + chunkSize - 1
			if end >= totalIterations {
				end = totalIterations - 1
			}
			wgSingle.Add(1)
			go func(start, end int) {
				defer wgSingle.Done()
				for iter := start; iter <= end; iter++ {
					i := iter/length + 1
					j := iter%length + 1
					a[i-1][j-1] = a[i-1][j-1] + 1
				}
			}(start, end)
		}
		wgSingle.Wait()
		//$omp end taskloop
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	fmt.Printf("a(50,50) =%3d\n", a[49][49])
}
