/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Data Race free matrix vector multiplication using target construct.
Each goroutine works on different elements of the result array.
*/

package main

import (
	"fmt"
	"sync"
)

const C = 100

func main() {
	// Allocate arrays
	a := make([]int, C)
	b := make([]int, C*C)
	c := make([]int, C)

	// Initialize arrays
	for i := 0; i < C; i++ {
		for j := 0; j < C; j++ {
			b[j+i*C] = 1
		}
		a[i] = 1
		c[i] = 0
	}

	// Simulate target teams distribute parallel for
	var wg sync.WaitGroup
	numGoroutines := 8
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(teamID int) {
			defer wg.Done()

			// Each team processes a chunk of rows
			start := teamID * C / numGoroutines
			end := (teamID + 1) * C / numGoroutines

			for i := start; i < end; i++ {
				// Matrix-vector multiplication for row i
				// No race: each goroutine writes to different c[i]
				for j := 0; j < C; j++ {
					c[i] += b[j+i*C] * a[j]
				}
			}
		}(t)
	}

	wg.Wait()

	// Check results
	for i := 0; i < C; i++ {
		if c[i] != C {
			fmt.Printf("Data Race\n")
			return
		}
	}

	fmt.Printf("Success: Matrix-vector multiplication completed correctly\n")
}
