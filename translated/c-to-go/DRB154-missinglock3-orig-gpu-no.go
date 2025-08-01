/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Concurrent accessing var may cause atomicity violation and inter region data race.
Lock and reduction clause avoids this. No Data Race Pair.
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var variable int64 = 0
	var wg sync.WaitGroup

	// Simulate teams distribute with reduction pattern
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of iterations
			start := teamID * 100 / numGoroutines
			end := (teamID + 1) * 100 / numGoroutines

			// Local accumulation (reduction pattern)
			localSum := int64(0)
			var localLock sync.Mutex

			for i := start; i < end; i++ {
				// Lock protects local operations
				localLock.Lock()
				localSum++
				localLock.Unlock()
			}

			// Atomic reduction to global variable
			atomic.AddInt64(&variable, localSum)
		}(t)
	}

	wg.Wait()
	fmt.Printf("%d \n", variable)
}
