/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Depend clause will ensure that there is no data race. There is an implicit barrier after tasks execution.
*/

package main

import (
	"fmt"
	"sync"
)

const C = 64

func main() {
	var a float32
	var x [C]float32
	var y [C]float32
	var wg sync.WaitGroup

	// Initialize arrays
	for i := 0; i < C; i++ {
		a = 5
		x[i] = 0
		y[i] = 3
	}

	// Simulate target with task dependencies
	// Each element has sequential dependency: multiply then add
	for i := 0; i < C; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Task 1: multiply (depends on x[index])
			localX := x[index]
			localX = a * localX

			// Task 2: add (depends on previous task completion)
			localX = localX + y[index]

			// Write back result
			x[index] = localX
		}(i)
	}

	wg.Wait()

	// Check results
	for i := 0; i < C; i++ {
		if x[i] != 3 {
			fmt.Printf("Data Race Detected\n")
			return
		}
	}

	fmt.Printf("Success: All computations completed correctly\n")
}
