//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Test if atomic can be recognized properly. No data races.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	var a int64 // Use int64 for atomic operations
	a = 0

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp atomic
			atomic.AddInt64(&a, 1) // No race - atomic operation
			//$omp end atomic
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("a = %8d\n", a)
}
