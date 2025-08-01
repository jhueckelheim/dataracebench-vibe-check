/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Vector addition followed by multiplication involving the same var should have a barrier in between.
Here we have an implicit barrier after parallel for regions. No data race pair.
*/

package main

import (
	"fmt"
	"sync"
)

const (
	N = 100
	C = 8
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

	// Simulate target parallel with implicit barriers
	for iteration := 0; iteration < N; iteration++ {
		var wg1 sync.WaitGroup
		var wg2 sync.WaitGroup

		// First parallel region: vector addition
		wg1.Add(C)
		for i := 0; i < C; i++ {
			go func(index int) {
				defer wg1.Done()
				temp[index] = b[index] + c[index]
			}(i)
		}
		wg1.Wait() // Implicit barrier after first parallel region

		// Second parallel region: multiplication
		wg2.Add(C)
		for i := C - 1; i >= 0; i-- {
			go func(index int) {
				defer wg2.Done()
				b[index] = temp[index] * a
			}(i)
		}
		wg2.Wait() // Implicit barrier after second parallel region
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
			fmt.Printf("expected %d real %d \n", val, b[i])
			return
		}
	}

	fmt.Printf("Success: Expected %d, all elements correct\n", val)
}
