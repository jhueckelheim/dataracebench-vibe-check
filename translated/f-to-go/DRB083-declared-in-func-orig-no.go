//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A variable is declared inside a function called within a parallel region.
//The variable should be private if it does not use static storage. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

func foo() {
	var q int // Local variable - private to each function call
	q = 0
	q = q + 1 // No race - each goroutine has its own copy of q
}

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			foo() // No race - local variables are private to each call
		}()
	}
	wg.Wait()
	//$omp end parallel
}
