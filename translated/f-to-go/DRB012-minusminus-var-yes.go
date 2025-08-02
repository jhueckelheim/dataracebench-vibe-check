//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The -1 operation is not protected, causing race condition.
//Data race pair: numNodes2@59:13:W vs. numNodes2@59:13:W

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var i, length, numNodes, numNodes2, argCount, allocStatus, ix int
	var rdErr error
	var args []string
	var x []int

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

	x = make([]int, length)

	numNodes = length
	numNodes2 = 0
	// initialize x()
	for i = 1; i <= length; i++ {
		if i%2 == 0 {
			x[i-1] = 5
		} else {
			x[i-1] = -5
		}
	}

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := numNodes / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := numNodes; start >= 1; start -= chunkSize {
		end := start - chunkSize + 1
		if end < 1 {
			end = 1
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i >= end; i-- {
				if x[i-1] <= 0 {
					numNodes2 = numNodes2 - 1 // Race condition - unprotected decrement
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("numNodes2 = %v\n", numNodes2)

	// deallocate(args,x)
}
