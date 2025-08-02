//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimensional array computation:
//default(none) to enforce explictly list all variables in data-sharing attribute clauses
//default(shared) to cover another option. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var length int
	var a, b [][]float64

	length = 100

	a = make([][]float64, length)
	b = make([][]float64, length)
	for i := 0; i < length; i++ {
		a[i] = make([]float64, length)
		b[i] = make([]float64, length)
	}

	//$omp parallel do default(none) shared(a) private(i,j)
	// default(none): explicitly specify all variables
	// shared(a): a is explicitly shared
	// private(i,j): i,j are explicitly private
	var wg1 sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 100 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= 100; start += chunkSize {
		end := start + chunkSize - 1
		if end > 100 {
			end = 100
		}
		wg1.Add(1)
		go func(start, end int) {
			defer wg1.Done()
			for i := start; i <= end; i++ { // i is private to each goroutine
				for j := 1; j <= 100; j++ { // j is private to each goroutine
					a[i-1][j-1] = a[i-1][j-1] + 1 // No race - proper partitioning
				}
			}
		}(start, end)
	}
	wg1.Wait()
	//$omp end parallel do

	//$omp parallel do default(shared) private(i,j)
	// default(shared): all variables shared unless explicitly private
	// private(i,j): i,j are explicitly private
	var wg2 sync.WaitGroup

	for start := 1; start <= 100; start += chunkSize {
		end := start + chunkSize - 1
		if end > 100 {
			end = 100
		}
		wg2.Add(1)
		go func(start, end int) {
			defer wg2.Done()
			for i := start; i <= end; i++ { // i is private to each goroutine
				for j := 1; j <= 100; j++ { // j is private to each goroutine
					b[i-1][j-1] = b[i-1][j-1] + 1 // No race - proper partitioning
				}
			}
		}(start, end)
	}
	wg2.Wait()
	//$omp end parallel do

	fmt.Printf("%f %f\n", a[49][49], b[49][49])

	// deallocate(a,b)
}
