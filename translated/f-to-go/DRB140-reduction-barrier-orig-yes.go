//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The assignment to a@21:9 is  not synchronized with the update of a@29:11 as a result of the
//reduction computation in the for loop.
//Data Race pair: a@21:9:W vs. a@24:30:W

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
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			//$omp master
			if threadID == 0 {
				a = 0 // RACE: Master sets a without barrier
			}
			//$omp end master
			// NO BARRIER HERE - causes race!

			//$omp do reduction(+:a)
			localA := 0
			chunkSize := 10 / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = 10
			}

			for i := start; i <= end; i++ {
				localA = localA + i // RACE: reduction on a while master might be setting a
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
