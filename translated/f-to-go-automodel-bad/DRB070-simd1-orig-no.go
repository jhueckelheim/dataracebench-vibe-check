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
	for i := 1; i <= length; i++ {
		a[i-1] = b[i-1] + c[i-1]
	}
	//$omp end simd
}
