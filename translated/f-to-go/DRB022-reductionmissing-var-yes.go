//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A kernel for two level parallelizable loop with reduction:
//if reduction(+:sum) is missing, there is race condition.
//Data race pairs:
//  getSum@60:13:W vs. getSum@60:13:W
//  getSum@60:13:W vs. getSum@60:22:R

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
	var getSum float32
	var rdErr error
	var args []string
	var u [][]float32

	length = 100
	getSum = 0.0

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

	u = make([][]float32, length)
	for i = 0; i < length; i++ {
		u[i] = make([]float32, length)
	}

	for i = 1; i <= length; i++ {
		for j = 1; j <= length; j++ {
			u[i-1][j-1] = 0.5
		}
	}

	//$omp parallel do private(temp, i, j)
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
			var temp float32 // private variable
			var i, j int     // private variables
			for i = start; i <= end; i++ {
				for j = 1; j <= length; j++ {
					temp = u[i-1][j-1]
					getSum = getSum + temp*temp // Race condition: missing reduction
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("sum = %v\n", getSum)

	// deallocate(args,u)
}
