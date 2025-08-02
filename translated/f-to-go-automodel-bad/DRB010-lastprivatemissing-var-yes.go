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
	var length, argCount, x int
	var args []string
	length = 10000

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
			for j := start; j <= end; j++ {
				x = j
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("x = %d\n", x)
} 