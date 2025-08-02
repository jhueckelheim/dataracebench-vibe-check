//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A single directive is used to protect a write. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var count int
	count = 0

	//$omp parallel shared(count)
	var wg sync.WaitGroup
	var once sync.Once
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			// Only one goroutine executes this block
			once.Do(func() {
				count = count + 1 // No race - only executed once
			})
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("count = %3d\n", count)
}
