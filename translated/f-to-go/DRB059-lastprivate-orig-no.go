//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Using lastprivate() to resolve an output dependence.
//
//Semantics of lastprivate (x):
//causes the corresponding original list item to be updated after the end of the region.
//The compiler/runtime copies the local value back to the shared one within the last iteration.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func foo() {
	var x int

	//$omp parallel do private(i) lastprivate(x)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 100 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	var mu sync.Mutex
	lastIteration := 100

	for start := 1; start <= 100; start += chunkSize {
		end := start + chunkSize - 1
		if end > 100 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localX := x // Copy of x for this goroutine
			for i := start; i <= end; i++ {
				localX = i
				// lastprivate semantics: if this is the last iteration, update shared x
				if i == lastIteration {
					mu.Lock()
					x = localX
					mu.Unlock()
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("x = %3d\n", x)
}

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			foo()
		}()
	}
	wg.Wait()
	//$omp end parallel
}
