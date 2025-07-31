/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
The barrier construct specifies an explicit barrier at the point at which the construct appears.
Barrier construct ensures that there is no data race.
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var variable int
	numThreads := runtime.NumCPU()

	// Barrier implementation using WaitGroup and channels
	var wg sync.WaitGroup
	var barrier1 sync.WaitGroup
	var barrier2 sync.WaitGroup

	// Set up barriers
	barrier1.Add(numThreads)
	barrier2.Add(numThreads)

	wg.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		go func(threadID int) {
			defer wg.Done()

			// First single execution (only one goroutine executes)
			if threadID == 0 {
				variable++
			}

			// Wait for all threads to reach barrier
			barrier1.Done()
			barrier1.Wait()

			// Second single execution (only one goroutine executes)
			if threadID == 0 {
				variable++
			}

			// Wait for all threads to reach second barrier
			barrier2.Done()
			barrier2.Wait()
		}(i)
	}

	wg.Wait()

	if variable != 2 {
		fmt.Printf("%d\n", variable)
	}

	// Return error if variable is not 2
	if variable != 2 {
		panic("Variable should be 2")
	}
}
