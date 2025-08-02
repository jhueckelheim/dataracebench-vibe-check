//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Reduction clause at line 23:34 will ensure there is no data race in var@27:13. No Data Race.

package main

import (
	"fmt"
	"sync"
)

func main() {
	var variable [8]int

	// Initialize
	for i := 1; i <= 8; i++ {
		variable[i-1] = 0
	}

	//$omp target map(tofrom:var) device(0)
	//$omp teams num_teams(1) thread_limit(1048)
	//$omp distribute parallel do reduction(+:var)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// Local copy for reduction
			localVar := [8]int{}
			
			//$omp simd
			for j := 1; j <= 8; j++ {
				localVar[j-1] = localVar[j-1] + 1 // No race - local reduction
			}
			//$omp end simd
			
			// Reduction operation
			mu.Lock()
			for j := 0; j < 8; j++ {
				variable[j] += localVar[j]
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	//$omp end distribute parallel do
	//$omp end teams
	//$omp end target

	for i := 1; i <= 8; i++ {
		if variable[i-1] != 20 {
			fmt.Printf("%d\n", variable[i-1])
		}
	}
}