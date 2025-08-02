//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The safelen(2) clause safelen(2)@23:16 guarantees that the vector code is safe for vectors up to 2 (inclusive).
//In the loop, m can be 2 or more for the correct execution. If the value of m is less than 2,
//the behavior is undefined. No Data Race in b[i]@25:9 assignment.

package main

import (
	"fmt"
)

func main() {
	var i, m, n int
	var b [4]float64

	m = 2
	n = 4

	//$omp simd safelen(2)
	for i = m + 1; i <= n; i++ {
		b[i-1] = b[i-m-1] - 1.0
	}
	//$omp end simd

	fmt.Printf("%f\n", b[2])
}
