//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This example is derived from an example by Simone Atzeni, NVIDIA.
//
//Description: Race on variable init if used master construct. The variable is written by the
//master thread and concurrently read by the others.
//
//Solution: master construct does not have an implicit barrier better
//use single at line 26. Fixed version for DRB124-master-orig-yes.c. No data race.

package main

import (
	"runtime"
	"sync"
)

func main() {
	var init, local int

	//$omp parallel shared(init) private(local)
	var wg sync.WaitGroup
	var once sync.Once
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			// Single directive HAS implicit barrier
			once.Do(func() {
				init = 10 // No race - only one execution
			})
			//$omp end single (implicit barrier here)

			local = init // No race - barrier ensures init is set
			_ = local    // Use local to avoid unused variable
		}()
	}
	wg.Wait()
	//$omp end parallel
}
