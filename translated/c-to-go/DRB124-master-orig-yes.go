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

Description: Race on variable init. The variable is written by the
master thread and concurrently read by the others.

Solution: master construct does not have an implicit barrier - better
use single. Data Race Pair: init (write vs. read)
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

	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			var local int

			// Master goroutine (equivalent to master directive without barrier)
			if threadID == 0 {
				init = 10 // Data race: master goroutine writes to init
			}

			// All goroutines (including master) read init immediately
			local = init // Data race: concurrent read while master may be writing
			_ = local
		}(t)
	}

	wg.Wait()
}
