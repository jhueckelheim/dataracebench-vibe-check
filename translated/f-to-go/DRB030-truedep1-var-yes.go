//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This program has data races due to true dependence within a loop.
//Data race pair: a[i+1]@51:9:W vs. a[i]@51:18:R

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

	for i := 1; i <= length; i++ {
		a[i-1] = i
	}

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := (length - 1) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= length-1; start += chunkSize {
		end := start + chunkSize - 1
		if end > length-1 {
			end = length - 1
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i] = a[i-1] + 1 // True dependence race condition
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("a(50)=%3d\n", a[49])

	// deallocate(args,a)
}
