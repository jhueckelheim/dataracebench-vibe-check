//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This one has data races due to true dependence.
//But data races happen at instruction level, not thread level.
//Data race pair: a[i+1]@32:9:W vs. a[i]@32:18:R

package main

import (
	"fmt"
)

func main() {
	var i, length int
	var a []int
	var b []int

	length = 100

	a = make([]int, length)
	b = make([]int, length)

	for i = 1; i <= length; i++ {
		a[i-1] = i
		b[i-1] = i + 1
	}

	//$omp simd
	for i = 1; i <= length-1; i++ {
		a[i] = a[i-1] + b[i-1]
	}

	for i = 1; i <= length; i++ {
		fmt.Printf("Values for i and a(i) are: %d %d\n", i, a[i-1])
	}
}
