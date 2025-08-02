//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The workshare construct is only available in Fortran. The workshare spreads work across the threads 
//executing the parallel. There is an implicit barrier. No data race.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var AA, BB, CC, res int

	BB = 1
	CC = 2

	//$omp parallel
	var wg sync.WaitGroup
	var barrier sync.WaitGroup
	numCPU := runtime.NumCPU()
	
	barrier.Add(numCPU) // Set up barrier between workshares

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			//$omp workshare
			// Workshare distributes work across threads
			if threadID == 0 { // Simulate work distribution
				AA = BB
				AA = AA + CC
			}
			//$omp end workshare (implicit barrier)
			
			barrier.Done()
			barrier.Wait() // Implicit barrier ensures first workshare completes

			//$omp workshare
			if threadID == 0 { // Second workshare
				res = AA * 2 // No race - barrier ensures AA is ready
			}
			//$omp end workshare
		}()
	}
	wg.Wait()
	//$omp end parallel

	if res != 6 {
		fmt.Printf("%d\n", res)
	}
}