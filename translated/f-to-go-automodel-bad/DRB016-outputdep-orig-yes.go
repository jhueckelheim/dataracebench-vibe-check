//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The loop in this example cannot be parallelized.
//
//This pattern has two pair of dependencies:
//1. loop carried output dependence
// x = .. :
//
//2. loop carried true dependence due to:
//.. = x;
// x = ..;
//Data race pairs: we allow two pairs to preserve the original code pattern.
// 1. x@48:16:R vs. x@49:9:W
// 2. x@49:9:W vs. x@49:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// globalArray module equivalent
var a []int

func useGlobalArray(length *int) {
	*length = 100
	a = make([]int, 100)
}

func main() {
	var length, x int

	length = 100
	x = 10

	useGlobalArray(&length)

	//$omp parallel do
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
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = x
				x = i
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("x = %d\n", x)
}
