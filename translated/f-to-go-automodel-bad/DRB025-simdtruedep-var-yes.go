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
	var i, length, argCount int
	var args []string
	var a []int
	var b []int

	length = 100

	argCount = len(os.Args) - 1
	if argCount == 0 {
		fmt.Println("No command line arguments provided.")
	}

	args = os.Args[1:]

	if argCount >= 1 {
		var rdErr error
		length, rdErr = strconv.Atoi(args[0])
		if rdErr != nil {
			fmt.Println("Error, invalid integer value.")
		}
	}

	a = make([]int, length)
	b = make([]int, length)

	for i = 1; i <= length; i++ {
		a[i-1] = i
		b[i-1] = i + 1
	}

	//$omp simd
	for i = 1; i <= length-1; i++ {
		a[i] = a[i-1] + b[i-1]
	}

	for i = 1; i <= length; i++ {
		fmt.Printf("Values for i and a(i) are: %d %d\n", i, a[i-1])
	}
}
