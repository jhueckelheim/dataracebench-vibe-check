//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Parallel region variable values are written to different output streams. No data race

package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func(tid int) {
			defer wg.Done()

			// Each thread writes to its own stream/output
			// In Go, we'll use different format strings to distinguish
			if tid == 0 {
				fmt.Fprintf(os.Stdout, "Hello World from thread %d\n", tid)
			} else if tid == 1 {
				fmt.Fprintf(os.Stderr, "Hello World from thread %d\n", tid)
			} else {
				fmt.Printf("Hello World from thread %d\n", tid) // Default output
			}
		}(threadID)
	}
	wg.Wait()
	//$omp end parallel
}
