//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two nested loops with loop-carried anti-dependence on the outer level.
//This is a variable-length array version in F95.
//Data race pair: a[i][j]@55:13:W vs. a[i+1][j]@55:31:R

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var i, j, length, argCount, allocStatus, ix int
	var rdErr error
	var args []string
	var a [][]float32
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

	a = make([][]float32, length)
	for i = 0; i < length; i++ {
		a[i] = make([]float32, length)
	}

	for i = 1; i <= length; i++ {
		for j = 1; j <= length; j++ {
			a[i-1][j-1] = 0.5
		}
	}

	//$omp parallel do private(j)
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
			var j int // private variable
			for i := start; i <= end; i++ {
				for j = 1; j <= length; j++ {
					a[i-1][j-1] = a[i-1][j-1] + a[i][j-1]
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("a(10,10) = %v\n", a[9][9])

	// deallocate(a)
	// deallocate(args)
}
