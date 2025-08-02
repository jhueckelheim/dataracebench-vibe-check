//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent access of var@30:13 without acquiring locks causes atomicity violation. Data race present.
//Data Race Pairs, var@30:13:W vs. var@30:22:R

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
	//$omp teams distribute parallel do
	// MISSING: reduction clause
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			//$omp simd
			for j := 1; j <= 16; j++ {
				variable[j-1] = variable[j-1] + 1 // RACE: Multiple threads + SIMD accessing shared variable
			}
			//$omp end simd
		}()
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", variable[15])
}