//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This one has race condition due to true dependence.
//But data races happen at instruction level, not thread level.
//Data race pair: a[i+1]@55:18:R vs. a[i]@55:9:W

package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var i, length, argCount, allocStatus, ix int
	var rdErr error
	var args []string
	var a, b []int

	length = 100

	argCount = len(os.Args) - 1
	if argCount == 0 {
		fmt.Printf("No command line arguments provided.\n")
	}

	args = make([]string, argCount)
	allocStatus = 0
	if allocStatus > 0 {
		fmt.Printf("Allocation error, program terminated.\n")
		os.Exit(1)
	}

	for ix = 1; ix <= argCount; ix++ {
		if ix <= len(os.Args)-1 {
			args[ix-1] = os.Args[ix]
		}
	}

	if argCount >= 1 {
		length, rdErr = strconv.Atoi(args[0])
		if rdErr != nil {
			fmt.Printf("Error, invalid integer value.\n")
		}
	}

	a = make([]int, length)
	b = make([]int, length)

	for i = 1; i <= length; i++ {
		a[i-1] = i
		b[i-1] = i + 1
	}

	//$omp simd
	// In Go, SIMD is handled by the compiler/runtime, we use regular loop
	for i = 1; i <= length-1; i++ {
		a[i] = a[i-1] + b[i-1] // True dependence: a[i+1] depends on a[i]
	}

	for i = 1; i <= length; i++ {
		fmt.Printf("Values for i and a(i) are: %d %d\n", i, a[i-1])
	}

	// deallocate(args,a,b)
}
