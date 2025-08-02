//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//omp_get_thread_num() is used to ensure serial semantics. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var numThreads int

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			// Simulate omp_get_thread_num() == 0 behavior
			if threadID == 0 {
				numThreads = numCPU // Equivalent to omp_get_num_threads()
			}
		}(i)
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("numThreads = %3d\n", numThreads)
}
