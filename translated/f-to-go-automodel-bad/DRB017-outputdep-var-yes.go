//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The loop in this example cannot be parallelized.
//
//Data race pairs: we allow two pairs to preserve the original code pattern.
// 1. x@50:16:R vs. x@51:9:W
// 2. x@51:9:W vs. x@51:9:W

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var length, x, argCount int
	var args []string
	var a []int

	length = 100
	x = 10

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

	a = make([]int, length)
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
				a[i-1] = x
				x = i
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("x= %3d  a(0)= %3d\n", x, a[0])
}
