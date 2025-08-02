//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent access of var@23 in an intra region. Lock ensures that there is no data race.

package main

import (
	"sync"
)

func main() {
	var lck sync.Mutex
	var variable int
	variable = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams num_teams(1)
	//$omp distribute parallel do
	// Single team ensures proper lock synchronization
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lck.Lock()   // Lock within single team
			variable = variable + 1 // No race - proper lock protection
			lck.Unlock()
		}()
	}
	wg.Wait()
	//$omp end distribute parallel do
	//$omp end teams
	//$omp end target
}