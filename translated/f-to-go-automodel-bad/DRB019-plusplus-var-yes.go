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
	var i, inLen, outLen, argCount int
	var args []string
	var input []int
	var output []int

	inLen = 1000
	outLen = 1

	argCount = len(os.Args) - 1
	if argCount == 0 {
		fmt.Println("No command line arguments provided.")
	}

	args = os.Args[1:]

	if argCount >= 1 {
		var rdErr error
		inLen, rdErr = strconv.Atoi(args[0])
		if rdErr != nil {
			fmt.Println("Error, invalid integer value.")
		}
	}

	input = make([]int, inLen)
	output = make([]int, inLen)

	for i = 1; i <= inLen; i++ {
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
				output[outLen-1] = input[i-1]
				outLen = outLen + 1
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("output(0)=%3d\n", output[0])
}
