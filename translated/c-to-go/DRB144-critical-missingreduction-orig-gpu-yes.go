/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
The increment is critical for the variable var, but the critical section is missing.
Therefore, there is a Data Race pair: var (write vs. write)
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
)

const N = 100

var variable int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup

	// Simulate target teams distribute parallel for
	numGoroutines := runtime.NumCPU() * 2 // Simulate many parallel threads like GPU
	wg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			defer wg.Done()

			// Each goroutine processes a chunk of iterations
			chunkSize := (N * 2) / numGoroutines
			start := goroutineID * chunkSize
			end := start + chunkSize
			if goroutineID == numGoroutines-1 {
				end = N * 2 // Handle remainder for last goroutine
			}

			for i := start; i < end; i++ {
				// Missing critical section - this should be protected but isn't
				variable++ // Data race: concurrent writes to variable without synchronization
			}
		}(g)
	}

	wg.Wait()

	fmt.Printf("%d\n", variable)
}
