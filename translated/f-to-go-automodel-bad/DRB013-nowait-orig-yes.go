//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This example is extracted from a paper:
//Ma etc. Symbolic Analysis of Concurrency Errors in OpenMP Programs, ICPP 2013
//
//Some threads may finish the for loop early and execute errors = dt[10]+1
//while another thread may still be simultaneously executing
//the for worksharing region by writing to d[9], causing data races.
//
//Data race pair: a[i]@41:21:R vs. a[10]@37:17:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, error, length, b int
	var a []int

	b = 5
	length = 1000

	a = make([]int, length)

	for i = 1; i <= length; i++ {
		a[i-1] = i
	}

	//$omp parallel shared(b, error)
	var wgOuter sync.WaitGroup
	wgOuter.Add(1)
	go func() {
		defer wgOuter.Done()
		//$omp parallel
		var wgInner sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := length / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		//$omp do
		for start := 1; start <= length; start += chunkSize {
			end := start + chunkSize - 1
			if end > length {
				end = length
			}
			wgInner.Add(1)
			go func(start, end int) {
				defer wgInner.Done()
				for i := start; i <= end; i++ {
					a[i-1] = b + a[i-1]*5
				}
			}(start, end)
		}
		//$omp end do nowait
		wgInner.Wait()

		//$omp single
		error = a[9] + 1
		//$omp end single
		//$omp end parallel
	}()
	wgOuter.Wait()
	//$omp end parallel

	fmt.Printf("error = %d\n", error)
}
