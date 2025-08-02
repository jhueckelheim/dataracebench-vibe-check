//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Example use of firstprivate(). No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable (replacing module)
var a []int

func foo(a []int, n, g int) {
	//$omp parallel do firstprivate(g)
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
		go func(start, end, gCopy int) { // gCopy is firstprivate equivalent
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = a[i-1] + gCopy // No race - proper firstprivate behavior
			}
		}(start, end, g)
	}
	wg.Wait()
	//$omp end parallel do
}

func main() {
	a = make([]int, 100)
	foo(a, 100, 7)
	fmt.Printf("%d\n", a[49])
}