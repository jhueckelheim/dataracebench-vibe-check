/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* This kernel imitates the nature of a program from the NAS Parallel Benchmarks 3.0 MG suit.
 * The private(i) and explicit barrier will ensure synchronized behavior.
 * No Data Race Pairs.
 */

package main

import (
	"fmt"
	"sync"
)

func main() {
	var q [10]float64
	var qq [10]float64
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var barrier sync.WaitGroup

	// Initialize arrays
	for i := 0; i < 10; i++ {
		qq[i] = float64(i)
		q[i] = float64(i)
	}

	// Simulate parallel default(shared)
	numGoroutines := 4
	wg.Add(numGoroutines)
	barrier.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Simulate for private(i)
			start := threadID * 10 / numGoroutines
			end := (threadID + 1) * 10 / numGoroutines

			for i := start; i < end; i++ {
				q[i] += qq[i]
			}

			// Critical section
			mutex.Lock()
			q[9] += 1.0
			mutex.Unlock()

			// Barrier
			barrier.Done()
			barrier.Wait()

			// Single section (only one thread executes)
			if threadID == 0 {
				q[9] = q[9] - 1.0
			}
		}(t)
	}

	wg.Wait()

	// Print results
	for i := 0; i < 10; i++ {
		fmt.Printf("%f %f\n", qq[i], q[i])
	}
}
