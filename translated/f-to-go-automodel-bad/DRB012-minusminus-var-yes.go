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
	var i, length, numNodes, numNodes2, argCount int
	var args []string
	var x []int

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

	x = make([]int, length)

	numNodes = length
	numNodes2 = 0
	//initialize x()
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
					numNodes2 = numNodes2 - 1
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("numNodes2 = %d\n", numNodes2)
}
