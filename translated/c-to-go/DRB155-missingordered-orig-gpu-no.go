/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
By utilizing the ordered construct the execution will be
sequentially consistent. No Data Race Pair.
*/

package main

import (
	"fmt"
)

const N = 100

func main() {
	var variable [N]int

	// Initialize array
	for i := 0; i < N; i++ {
		variable[i] = 0
	}

	// Simulate ordered parallel execution - must maintain sequential order
	// OpenMP ordered ensures sequential execution despite parallel context
	// In Go, this means actually executing sequentially for correctness
	for i := 1; i < N; i++ {
		variable[i] = variable[i-1] + 1
	}

	// Check results
	for i := 0; i < N; i++ {
		if variable[i] != i {
			fmt.Printf("Data Race Present")
			return
		}
	}

	fmt.Printf("Success: All elements correctly ordered\n")
}
