//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//To avoid data race, the initialization of the original list item "a" should complete before any
//update of a as a result of the reduction clause. This can be achieved by adding an explicit
//barrier after the assignment a=0@22:9, or by enclosing the assignment a=0@22:9 in a single directive
//or by initializing a@21:7 before the start of the parallel region. No data race pair

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var a int

	//$omp parallel shared(a) private(i)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var barrier sync.WaitGroup
	numCPU := runtime.NumCPU()

	barrier.Add(numCPU) // Set up barrier

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			//$omp master
			if threadID == 0 {
				a = 0 // Master sets a
			}
			//$omp end master

			//$omp barrier
			barrier.Done()
			barrier.Wait() // Explicit barrier ensures a=0 completes first!

			//$omp do reduction(+:a)
			localA := 0
			chunkSize := 10 / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = 10
			}

			for i := start; i <= end; i++ {
				localA = localA + i // No race - barrier ensures proper ordering
			}

			// Reduction
			mu.Lock()
			a += localA
			mu.Unlock()
			//$omp end do

			//$omp single
			if threadID == 0 {
				fmt.Printf("Sum is %d\n", a)
			}
			//$omp end single
		}(threadID)
	}
	wg.Wait()
	//$omp end parallel
}
