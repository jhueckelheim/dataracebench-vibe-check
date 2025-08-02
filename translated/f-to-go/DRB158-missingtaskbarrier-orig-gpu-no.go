//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Depend clause at line 29 and 33 will ensure that there is no data race.

package main

import (
	"fmt"
	"sync"
)

// Package-level variables (module equivalent)
var a int
var x, y [64]int

func main() {
	// Initialize arrays
	for i := 1; i <= 64; i++ {
		x[i-1] = 0
		y[i-1] = 3
	}

	a = 5

	//$omp target map(to:y,a) map(tofrom:x) device(0)
	var allTasksWg sync.WaitGroup
	
	for i := 1; i <= 64; i++ {
		var taskWg sync.WaitGroup
		
		//$omp task depend(inout:x(i))
		taskWg.Add(1)
		allTasksWg.Add(1)
		go func(i int) {
			defer taskWg.Done()
			defer allTasksWg.Done()
			x[i-1] = a * x[i-1] // First task on x[i]
		}(i)

		//$omp task depend(inout:x(i))
		allTasksWg.Add(1)
		go func(i int) {
			defer allTasksWg.Done()
			taskWg.Wait() // Wait for dependency on x[i]
			x[i-1] = x[i-1] + y[i-1] // Second task depends on first
		}(i)
	}
	//$omp end target

	for i := 1; i <= 64; i++ {
		if x[i-1] != 3 {
			fmt.Printf("%d\n", x[i-1])
		}
	}

	//$omp taskwait
	allTasksWg.Wait()
}