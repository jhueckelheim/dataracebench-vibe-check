//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//use of omp target: len is not mapped. It should be firstprivate within target. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func main() {
	var length int
	var a []int

	length = 100 // Initialize length (was uninitialized in Fortran - undefined behavior)
	a = make([]int, length)

	// Initialize array
	for i := 1; i <= length; i++ {
		a[i-1] = i
	}

	//$omp target map(a(1:len))
	//$omp parallel do
	// In Go, we simulate target offloading with regular parallel execution
	// len is implicitly firstprivate (captured by value in goroutines)
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
		go func(start, end, lengthCopy int) { // lengthCopy is firstprivate equivalent
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = a[i-1] + 1 // No race - proper data partitioning
			}
		}(start, end, length)
	}
	wg.Wait()
	//$omp end target

	// deallocate(a) - handled automatically by Go's garbage collector
}
