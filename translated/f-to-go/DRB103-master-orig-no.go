//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A master directive is used to protect memory accesses. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var k int

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			//$omp master
			// Only master thread (thread 0) executes this block
			if threadID == 0 {
				k = numCPU // Equivalent to omp_get_num_threads()
				fmt.Printf("Number of threads requested = %8d\n", k)
			}
			//$omp end master
		}(threadID)
	}
	wg.Wait()
	//$omp end parallel
}
