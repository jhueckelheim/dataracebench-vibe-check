/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* This kernel imitates the nature of a program from the NAS Parallel Benchmarks 3.0 MG suit.
 * There is no data race pairs, example of a threadprivate var and update by TID==0 only.
 */

package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var j, k float64
	var jMutex, kMutex sync.Mutex

	// Simulate parallel for default(shared)
	numGoroutines := 4
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Thread-private x array (each goroutine has its own)
			var x [20]float64

			// Each thread processes a chunk of i values
			start := threadID * 20 / numGoroutines
			end := (threadID + 1) * 20 / numGoroutines

			for i := start; i < end; i++ {
				x[i] = -1.0

				// Only thread 0 updates shared variables
				if threadID == 0 {
					jMutex.Lock()
					j = x[0]
					jMutex.Unlock()
				}

				if threadID == 0 {
					kMutex.Lock()
					k = float64(i) + 0.05
					kMutex.Unlock()
				}
			}
		}(t)
	}

	wg.Wait()

	fmt.Printf("%f %f\n", j, k)
}
