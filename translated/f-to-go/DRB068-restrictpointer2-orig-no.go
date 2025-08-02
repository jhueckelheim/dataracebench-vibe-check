//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//micro-bench equivalent to the restrict keyword in C-99 in F95. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func foo(n int) ([]int, []int, []int, []int) {
	// Allocate separate arrays (no aliasing like restrict pointers)
	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n)
	d := make([]int, n)

	// Initialize arrays
	for i := 1; i <= n; i++ {
		b[i-1] = i
		c[i-1] = i
	}

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := n / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= n; start += chunkSize {
		end := start + chunkSize - 1
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = b[i-1] + c[i-1] // No race - a, b, c are separate arrays (restrict-like)
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	if a[499] != 1000 {
		fmt.Printf("%d\n", a[499])
	}

	return a, b, c, d
}

func main() {
	n := 1000

	// Create initial arrays
	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n)
	d := make([]int, n)

	a, b, c, d = foo(n)

	// In Go, slices are automatically garbage collected
	_ = a
	_ = b
	_ = c
	_ = d
}
