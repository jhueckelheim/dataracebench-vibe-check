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

// Module DRB068 translated to package-level functions
func foo(n int, a, b, c, d []int) {
	a = make([]int, n)
	b = make([]int, n)
	c = make([]int, n)
	d = make([]int, n)

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
				a[i-1] = b[i-1] + c[i-1]
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	if a[499] != 1000 {
		fmt.Printf("%d\n", a[499])
	}
}

func main() {
	var n int = 1000
	var a, b, c, d []int

	a = make([]int, n)
	b = make([]int, n)
	c = make([]int, n)
	d = make([]int, n)

	foo(n, a, b, c, d)
}
