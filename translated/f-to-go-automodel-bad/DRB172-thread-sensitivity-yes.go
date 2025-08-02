//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Iteration 1 and 2 can have conflicting writes to a(1). But if they are scheduled to be run by
// the same thread, dynamic tools may miss this.
// Data Race Pair, a(0)@39:9:W vs. a(i)@40:22:W

package main

import (
	"runtime"
	"sync"
)

func loadFromInput(a []int, N int) {
	for i := 1; i <= N; i++ {
		a[i-1] = i
	}
}

func main() {
	var i, N, argCount, allocStatus, rdErr, ix int
	var a []int

	N = 100

	a = make([]int, N)

	loadFromInput(a, N)

	//$omp parallel do shared(a)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := N / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= N; start += chunkSize {
		end := start + chunkSize - 1
		if end > N {
			end = N
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = i
				if i == 2 {
					a[0] = 1
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}
