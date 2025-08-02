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
	"sync"
)

// Module DRB105 translated to package-level variables and functions
var input int

func fib(n int) int {
	var i, j, r int

	if n < 2 {
		r = n
	} else {
		var wg sync.WaitGroup
		//$omp task shared(i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			i = fib(n - 1)
		}()
		//$omp end task
		//$omp task shared(j)
		wg.Add(1)
		go func() {
			defer wg.Done()
			j = fib(n - 2)
		}()
		//$omp end task
		//$omp taskwait
		wg.Wait()
		r = i + j
	}
	return r
}

func main() {
	var result int
	input = 30

	//$omp parallel
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp single
		result = fib(input)
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	fmt.Printf("Fib for %8d =%8d\n", input, result)
}
