//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Test if the semantics of omp_get_thread_num() is correctly recognized.
//Thread with id 0 writes numThreads while other threads read it, causing data races.
//Data race pair: numThreads@22:9:W vs. numThreads@24:31:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var numThreads int
	numThreads = 0

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			// Simulate omp_get_thread_num() behavior
			if threadID == 0 {
				numThreads = numCPU // RACE: Write to shared numThreads
			} else {
				fmt.Printf("numThreads = %d\n", numThreads) // RACE: Read shared numThreads
			}
		}(threadID)
	}
	wg.Wait()
	//$omp endparallel
}
