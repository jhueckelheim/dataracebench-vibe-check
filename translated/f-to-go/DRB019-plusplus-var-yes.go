//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Race condition on outLen due to unprotected writes.
//Adding private (outLen) can avoid race condition. But it is wrong semantically.
//
//Data race pairs: we allow two pair to preserve the original code pattern.
//1. outLen@60:9:W vs. outLen@60:9:W
//2. output[]@59:9:W vs. output[]@59:9:W

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var inLen, outLen, argCount, allocStatus, ix int
	var rdErr error
	var args []string
	var input, output []int

	inLen = 1000
	outLen = 1

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
		inLen, rdErr = strconv.Atoi(args[0])
		if rdErr != nil {
			fmt.Printf("Error, invalid integer value.\n")
		}
	}

	input = make([]int, inLen)
	output = make([]int, inLen)

	for i := 1; i <= inLen; i++ {
		input[i-1] = i
	}

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := inLen / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= inLen; start += chunkSize {
		end := start + chunkSize - 1
		if end > inLen {
			end = inLen
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				output[outLen-1] = input[i-1] // Race condition: accessing outLen
				outLen = outLen + 1           // Race condition: modifying outLen
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	// Handle output(0) - in Fortran this would be out of bounds access
	var output0 int
	if len(output) > 0 {
		output0 = output[0] // This represents output(1) in Fortran
	}
	fmt.Printf("output(0)=%3d\n", output0)

	// deallocate(input,output,args)
}
