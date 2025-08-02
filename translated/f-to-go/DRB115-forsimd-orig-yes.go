//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This one has data races due to true dependence.
//But data races happen at both instruction and thread level.
//Data race pair: a[i+1]@31:9:W vs. a[i]@31:16:R

package main

import (
	"fmt"
)

func main() {
	var length int
	var a, b []int

	length = 100
	a = make([]int, length)
	b = make([]int, length)

	// Initialize arrays
	for i := 1; i <= length; i++ {
		a[i-1] = i
		b[i-1] = i + 1
	}

	//$omp simd
	// SIMD vectorization with true dependence creates race at instruction level
	for i := 1; i <= length-1; i++ {
		a[i] = a[i-1] + b[i-1] // RACE: True dependence in vectorized execution
		// a[i+1] (write) vs a[i] (read) creates dependency
		// When vectorized, multiple iterations execute simultaneously
	}

	fmt.Printf("a(50) = %d\n", a[49])

	// deallocate(a,b)
}
