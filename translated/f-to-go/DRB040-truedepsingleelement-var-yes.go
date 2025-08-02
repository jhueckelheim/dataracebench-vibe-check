//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Data race pair: a[i]@49:9:W vs. a[0]@49:16:R

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var length, argCount, allocStatus, ix int
	var rdErr error
	var args []string
	var a []int

	length = 1000

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

	a[0] = 2

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := length / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = a[i-1] + a[0] // Race condition: all elements access a[1] (a[0] in Go)
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	// Handle a(0) access - in Fortran this would be out of bounds
	var a0 int
	if length > 0 {
		a0 = a[0] // Fortran a(0) is out of bounds, but we'll use a[0]
	}
	fmt.Printf("a(0) = %3d\n", a0)

	// deallocate(args,a)
}