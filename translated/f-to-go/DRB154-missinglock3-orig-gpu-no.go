//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent accessing var@25:9 may cause atomicity violation and inter region data race.
//Lock and reduction clause at line 22, avoids this. No Data Race Pair.

package main

import (
	"fmt"
	"sync"
)

func main() {
	var lck sync.Mutex
	var variable int
	variable = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute reduction(+:var)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numTeams := 4

	for team := 0; team < numTeams; team++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localVar := 0
			for i := 1; i <= 100; i++ {
				lck.Lock()   // Lock protection
				localVar = localVar + 1 // No race - reduction + lock
				lck.Unlock()
			}
			// Reduction
			mu.Lock()
			variable += localVar
			mu.Unlock()
		}()
	}
	wg.Wait()
	//$omp end teams distribute
	//$omp end target

	fmt.Printf("%d\n", variable)
}