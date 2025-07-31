/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
This example is referred from DRACC by Adrian Schmitz et al.
Vector addition followed by multiplication involving the same var should have a barrier in between.
omp distribute directive does not have implicit barrier. This will cause data race.
Data Race Pair: b[i]:R vs. b[i]:W
*/

package main

import (
	"fmt"
	"sync"
)

const (
	N = 100
	C = 16
)

func main() {
	var a int
	var b [C]int
	var c [C]int
	var temp [C]int

	// Initialize arrays
	for i := 0; i < C; i++ {
		b[i] = 0
		c[i] = 2
		temp[i] = 0
	}
	a = 2

	// Simulate teams with distribute (NO implicit barriers)
	for iteration := 0; iteration < N; iteration++ {
		var wg sync.WaitGroup

		// Launch both distribute operations concurrently without barrier
		// This creates the race condition
		wg.Add(2)

		// First distribute: vector addition
		go func() {
			defer wg.Done()
			for i := 0; i < C; i++ {
				temp[i] = b[i] + c[i] // RACE: Reading b[i] while second goroutine writes to it
			}
		}()

		// Second distribute: multiplication (runs concurrently!)
		go func() {
			defer wg.Done()
			for i := C - 1; i >= 0; i-- {
				b[i] = temp[i] * a // RACE: Writing b[i] while first goroutine reads it
			}
		}()

		wg.Wait()
	}

	// Calculate expected value
	val := 0
	for i := 0; i < N; i++ {
		val = val + 2
		val = val * 2
	}

	// Check results
	for i := 0; i < C; i++ {
		if b[i] != val {
			fmt.Printf("index: %d val: %d\n", i, b[i])
		}
	}
}
