//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This is a program based on a test contributed by Yizi Gu@Rice Univ.
//Proper user of ordered directive and clause, no data races

package main

import (
	"fmt"
)

func main() {
	var x int
	x = 0

	//$omp parallel do ordered
	// With proper ordered directive - sequential execution maintains correctness
	for i := 1; i <= 100; i++ {
		//$omp ordered
		// Sequential execution ensures no race
		x = x + 1 // No race - proper ordering ensures sequential access
		//$omp end ordered
	}
	//$omp end parallel do

	fmt.Printf("x = %d\n", x)
}
