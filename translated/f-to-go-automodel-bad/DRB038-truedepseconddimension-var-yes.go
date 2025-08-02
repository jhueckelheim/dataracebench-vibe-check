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
	var i, j, n, m, length, argCount int
	var args []string
	var b [][]float64

	length = 1000

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

	n = length
	m = length

	b = make([][]float64, n)
	for i := range b {
		b[i] = make([]float64, m)
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
			go func(i, start, end int) {
				defer wg.Done()
				for j := start; j <= end; j++ {
					b[i-1][j-1] = b[i-1][j-2]
				}
			}(i, start, end)
		}
		wg.Wait()
		//$omp end parallel do
	}

	// Uncomment to print a value for checking
	// fmt.Printf("b(5,5) =%20.6f\n", b[4][4])
}
