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
			for i := start; i <= end; i++ {
				for j := 2; j <= m; j++ {
					b[i-1][j-1] = b[i-2][j-2]
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("b(500,500) =%10.6f\n", b[499][499])
}
