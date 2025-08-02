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
	var mu sync.Mutex
	numCPU := runtime.NumCPU()
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//$omp single
			mu.Lock()
			init = 10
			mu.Unlock()
			//$omp end single
			local = init
		}()
	}
	wg.Wait()
	//$omp end parallel
}
