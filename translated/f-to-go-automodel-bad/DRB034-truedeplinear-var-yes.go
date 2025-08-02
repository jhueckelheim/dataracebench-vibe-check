//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A linear expression is used as array subscription.
//Data race pair: a[2*i+1]@53:9:W vs. a[i]@53:18:R

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var i, length, uLen, argCount int
	var args []string
	var a []int

	length = 2000

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

	for i = 1; i <= length; i++ {
		a[i-1] = i
	}

	uLen = length / 2

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := uLen / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= uLen; start += chunkSize {
		end := start + chunkSize - 1
		if end > uLen {
			end = uLen
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[2*i-1] = a[i-1] + 1
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}
