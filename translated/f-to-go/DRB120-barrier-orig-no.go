//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The barrier construct specifies an explicit barrier at the point at which the construct appears.
//Barrier construct at line:27 ensures that there is no data race.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var variable int

	//$omp parallel shared(var)
	var wg sync.WaitGroup
	var once1, once2 sync.Once
	var barrier sync.WaitGroup
	numCPU := runtime.NumCPU()

	// Set up barrier for all threads
	barrier.Add(numCPU)

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			once1.Do(func() {
				variable = variable + 1 // No race - only one execution
			})
			//$omp end single

			//$omp barrier
			barrier.Done()
			barrier.Wait() // All threads wait here

			//$omp single
			once2.Do(func() {
				variable = variable + 1 // No race - barrier ensures proper ordering
			})
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("var = %3d\n", variable)
}
