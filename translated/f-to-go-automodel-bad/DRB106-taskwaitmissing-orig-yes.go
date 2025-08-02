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
	"sync"
)

// Module DRB106 translated to package-level variables and functions
var input int32

func fib(n int32) int32 {
	var i, j, r int32

	if n < 2 {
		r = n
	} else {
		//$omp task shared(i)
		go func() {
			i = fib(n - 1)
		}()
		//$omp end task
		//$omp task shared(j)
		go func() {
			j = fib(n - 2)
		}()
		//$omp end task
		r = i + j
	}
	//$omp taskwait
	return r
}

func main() {
	var result int32
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

	fmt.Printf("Fib for %d = %d\n", input, result)
}
