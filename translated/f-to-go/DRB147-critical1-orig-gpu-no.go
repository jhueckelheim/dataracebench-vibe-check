//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent access on same variable var@23 and var@25 leads to the race condition if two different
//locks are used. This is the reason here we have used the atomic directive to ensure that addition
//and subtraction are not interleaved. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	var variable int64 // Use int64 for atomic operations
	variable = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()

	for team := 0; team < numTeams; team++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i <= 100; i++ {
				//$omp atomic
				atomic.AddInt64(&variable, 1) // No race - atomic operation
				//$omp atomic
				atomic.AddInt64(&variable, -2) // No race - atomic operation
			}
		}()
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", variable)
}