/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
This example is derived from an example by Simone Atzeni, NVIDIA.

Description: Fixed version for DRB124-master-orig-yes.c. No data race.
The single directive has an implicit barrier, ensuring all threads wait
before reading the variable.
*/
package main

import (
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var init int
	numThreads := runtime.NumCPU()
	var wg sync.WaitGroup
	var barrier sync.WaitGroup

	// Set up barrier for all threads
	barrier.Add(numThreads)

	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			var local int

			// Single execution (equivalent to single directive with barrier)
			if threadID == 0 {
				init = 10 // Only one goroutine writes
			}

			// Implicit barrier - wait for single execution to complete
			barrier.Done()
			barrier.Wait()

			// All goroutines read after barrier
			local = init // No data race: read after write is complete
			_ = local
		}(t)
	}

	wg.Wait()
}
