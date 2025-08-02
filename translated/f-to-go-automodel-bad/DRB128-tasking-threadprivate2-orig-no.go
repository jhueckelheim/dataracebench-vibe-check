//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The scheduling constraints prohibit a thread in the team from executing
//a new task that modifies tp while another such task region tied to
//the same thread is suspended. Therefore, the value written will
//persist across the task scheduling point.
//No Data Race

package main

import (
	"fmt"
)

// Module DRB128 translated to package-level variables
var tp, var1 int

func foo() {
	//$omp task
	go func() {
		//$omp task
		go func() {
			tp = 1
			//$omp task
			//$omp end task
			var1 = tp
		}()
		//$omp end task
	}()
	//$omp end task
}

func main() {
	foo()
	fmt.Printf("%d\n", var1)
}
