//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//One dimension array computation with a vectorization directive. No data race pairs.

package main

func main() {
	var length int
	var a, b, c []int

	length = 100
	a = make([]int, length)
	b = make([]int, length)
	c = make([]int, length)

	//$omp simd
	// In Go, we rely on the compiler for potential vectorization
	for i := 1; i <= length; i++ {
		a[i-1] = b[i-1] + c[i-1] // No race - sequential execution with potential vectorization
	}
	//$omp end simd

	// deallocate(a,b,c) - handled automatically by Go's garbage collector
}
