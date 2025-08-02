//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent access of var@26:13 has no atomicity violation. No data race present.

package main

import (
	"fmt"
	"sync"
)

// Package-level variables (module equivalent)
var variable [16]int

func main() {
	// Initialize
	for i := 1; i <= 16; i++ {
		variable[i-1] = 0
	}

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do reduction(+:var)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// Local copy for reduction
			localVar := [16]int{}
			
			//$omp simd
			for j := 1; j <= 16; j++ {
				localVar[j-1] = localVar[j-1] + 1 // No race - reduction protects
			}
			//$omp end simd
			
			// Reduction operation
			mu.Lock()
			for j := 0; j < 16; j++ {
				variable[j] += localVar[j]
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	for i := 1; i <= 16; i++ {
		if variable[i-1] != 20 {
			fmt.Printf("%d %d\n", variable[i-1], i)
		}
	}
}