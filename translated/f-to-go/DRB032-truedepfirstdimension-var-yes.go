//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The outer loop has a loop-carried true dependence.
//Data race pair: b[i][j]@56:13:W vs. b[i-1][j-1]@56:22:R

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var i, j, n, m, length, argCount, allocStatus, ix int
	var rdErr error
	var args []string
	var b [][]float32

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

	n = length
	m = length
	b = make([][]float32, n)
	for i = 0; i < n; i++ {
		b[i] = make([]float32, m)
	}

	for i = 1; i <= n; i++ {
		for j = 1; j <= m; j++ {
			b[i-1][j-1] = 0.5
		}
	}

	//$omp parallel do private(j)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := (n - 1) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 2; start <= n; start += chunkSize {
		end := start + chunkSize - 1
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			var j int // private variable
			for i := start; i <= end; i++ {
				for j = 2; j <= m; j++ {
					b[i-1][j-1] = b[i-2][j-2] // True dependence race condition
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("b(500,500) = %10.6f\n", b[499][499])

	// deallocate(args,b)
}
