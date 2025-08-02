//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimension array computation with a vetorization directive
//collapse(2) makes simd associate with 2 loops.
//Loop iteration variables should be predetermined as lastprivate. No data race pairs.

package main

import (
	"fmt"
)

func main() {
	var a, b, c [][]float64
	var length int

	length = 100
	a = make([][]float64, length)
	b = make([][]float64, length)
	c = make([][]float64, length)

	for i := range a {
		a[i] = make([]float64, length)
		b[i] = make([]float64, length)
		c[i] = make([]float64, length)
	}

	for i := 1; i <= length; i++ {
		for j := 1; j <= length; j++ {
			a[i-1][j-1] = float64(i) / 2.0
			b[i-1][j-1] = float64(i) / 3.0
			c[i-1][j-1] = float64(i) / 7.0
		}
	}

	//$omp simd collapse(2)
	for i := 1; i <= length; i++ {
		for j := 1; j <= length; j++ {
			c[i-1][j-1] = a[i-1][j-1] * b[i-1][j-1]
		}
	}
	//$omp end simd

	fmt.Printf("c(50,50) = %f\n", c[49][49])
}
