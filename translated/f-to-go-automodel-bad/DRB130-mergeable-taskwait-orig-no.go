//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Taken from OpenMP Examples 5.0, example tasking.12.c
//x is a shared variable the outcome does not depend on whether or not the task is merged (that is,
//the task will always increment the same variable and will always compute the same value for x).

package main

import (
	"fmt"
)

func main() {
	var x int
	x = 2

	//$omp task shared(x) mergeable
	go func() {
		x = x + 1
	}()
	//$omp end task

	fmt.Printf("x =%8d\n", x)
}
