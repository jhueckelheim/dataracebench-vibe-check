//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Only the outmost loop can be parallelized in this program.
//Data race pair: b[i][j]@51:13:W vs. b[i][j-1]@51:22:R

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var i, n, m, length, argCount, allocStatus, ix int
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

	b = make([][]float32, length)
	for i = 0; i < length; i++ {
		b[i] = make([]float32, length)
	}

	for i = 1; i <= n; i++ {
		//$omp parallel do
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := (m - 1) / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 2; start <= m; start += chunkSize {
			end := start + chunkSize - 1
			if end > m {
				end = m
			}
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for j := start; j <= end; j++ {
					b[i-1][j-1] = b[i-1][j-2] // True dependence on second dimension
				}
			}(start, end)
		}
		wg.Wait()
		//$omp end parallel do
	}
	// Commented out print statement preserved from original
	// fmt.Printf("b(5,5) = %20.6f\n", b[4][4])

	// deallocate(args,b)
}