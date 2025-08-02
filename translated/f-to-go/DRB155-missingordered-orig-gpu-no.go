//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//By utilizing the ordered construct @23 the execution will be sequentially consistent.
//No Data Race Pair.

package main

import (
	"fmt"
)

func main() {
	var variable [100]int

	// Initialize
	for i := 1; i <= 100; i++ {
		variable[i-1] = 1
	}

	//$omp target map(tofrom:var) device(0)
	//$omp parallel do ordered
	// Ordered ensures sequential consistency
	for i := 2; i <= 100; i++ {
		//$omp ordered
		// Sequential execution maintains proper ordering
		variable[i-1] = variable[i-2] + 1 // No race - sequential execution
		//$omp end ordered
	}
	//$omp end parallel do
	//$omp end target

	for i := 1; i <= 100; i++ {
		if variable[i-1] != i {
			fmt.Printf("Data Race Present\n")
		}
	}
}