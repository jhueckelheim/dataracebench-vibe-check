//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//* This is a program based on a test contributed by Yizi Gu@Rice Univ.
//* Classic Fibonacci calculation using task+taskwait. No data races.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable (module equivalent)
var input int

func fib(n int) int {
	if n < 2 {
		return n
	} else {
		var i, j int
		var wg sync.WaitGroup

		//$omp task shared(i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			i = fib(n - 1) // No race - taskwait ensures proper synchronization
		}()

		//$omp task shared(j)
		wg.Add(1)
		go func() {
			defer wg.Done()
			j = fib(n - 2) // No race - taskwait ensures proper synchronization
		}()

		//$omp taskwait
		wg.Wait() // Wait for both tasks to complete

		return i + j
	}
}

func main() {
	var result int
	input = 30

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			if threadID == 0 { // Only one thread executes single
				result = fib(input)
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("Fib for %8d = %8d\n", input, result)
}
