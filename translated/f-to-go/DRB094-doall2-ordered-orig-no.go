//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimensional array computation:
//ordered(2) is used to associate two loops with omp for.
//The corresponding loop iteration variables are private.
//
//ordered(n) is an OpenMP 4.5 addition. No data race pairs.

package main

import (
	"fmt"
)

// Package-level variable (module equivalent)
var a [][]int

func main() {
	var length int
	length = 100

	a = make([][]int, length)
	for i := 0; i < length; i++ {
		a[i] = make([]int, length)
	}

	//$omp parallel do ordered(2)
	// ordered(2) ensures ordered execution based on dependency constraints
	// In Go, we simulate this with sequential execution for correctness
	for i := 1; i <= length; i++ {
		for j := 1; j <= length; j++ {
			a[i-1][j-1] = a[i-1][j-1] + 1

			//$omp ordered depend(sink:i-1,j) depend(sink:i,j-1)
			// Dependencies ensure proper ordering - simulated with sequential execution
			fmt.Printf("test i = %d  j = %d\n", i, j)
			//$omp ordered depend(source)
		}
	}
	//$omp end parallel do
}
