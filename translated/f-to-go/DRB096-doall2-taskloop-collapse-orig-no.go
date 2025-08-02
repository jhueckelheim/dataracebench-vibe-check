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

// Package-level variable (module equivalent)
var a [][]int

func main() {
	var length int
	length = 100

	a = make([][]int, length)
	for i := 0; i < length; i++ {
		a[i] = make([]int, length)
	}

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			if threadID == 0 { // Only one thread executes single region
				//$omp taskloop collapse(2)
				// collapse(2) flattens nested loops for task distribution
				var taskWg sync.WaitGroup
				totalIterations := length * length

				for iteration := 0; iteration < totalIterations; iteration++ {
					taskWg.Add(1)
					go func(iteration int) { // iteration variables are private to each task
						defer taskWg.Done()
						// Convert flat iteration back to 2D indices
						i := (iteration / length) + 1 // Fortran 1-based
						j := (iteration % length) + 1 // Fortran 1-based
						a[i-1][j-1] = a[i-1][j-1] + 1 // No race - proper task partitioning
					}(iteration)
				}
				taskWg.Wait()
				//$omp end taskloop
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("a(50,50) = %3d\n", a[49][49])
}
