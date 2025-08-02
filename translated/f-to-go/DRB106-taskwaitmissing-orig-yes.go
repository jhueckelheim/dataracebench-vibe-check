//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//* This is a program based on a test contributed by Yizi Gu@Rice Univ.
// * Classic Fibonacci calculation using task but missing taskwait.
// * Data races pairs: i@29:13:W vs. i@34:17:R
// *                   j@32:13:W vs. j@34:19:R

//check on the unsgined part

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable (module equivalent)
var input int32

func fib(n int32) int32 {
	if n < 2 {
		return n
	} else {
		var i, j int32

		//$omp task shared(i)
		go func() {
			i = fib(n - 1) // RACE: No synchronization before reading i
		}()

		//$omp task shared(j)
		go func() {
			j = fib(n - 2) // RACE: No synchronization before reading j
		}()

		// MISSING: !$omp taskwait - this causes the race condition
		return i + j // RACE: Reading i and j before tasks complete
	}
	//$omp taskwait (misplaced - after return, never reached)
}

func main() {
	var result int32
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

	fmt.Printf("Fib for %d = %d\n", input, result)
}
