//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This loop has loop-carried output-dependence due to x=... at line 44.
//The problem can be solved by using lastprivate(x) .
//Data race pair: x@44:9:W vs. x@44:9:W

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var length, argCount, allocStatus, ix, x int
	var rdErr error
	var args []string
	length = 10000

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

	//$omp parallel do private(i)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := (length + 1) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 0; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			var i int // private variable
			for i = start; i <= end; i++ {
				x = i // This creates the data race - x is shared, not lastprivate
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("x = %v\n", x)

	// deallocate(args)
}
