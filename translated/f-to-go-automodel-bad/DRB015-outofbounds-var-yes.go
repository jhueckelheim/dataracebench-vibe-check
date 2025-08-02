//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The outmost loop is be parallelized.
//But the inner level loop has out of bound access for b[i][j] when i equals to 1.
//This will case memory access of a previous column's last element.
//
//For example, an array of 4x4:
//    j=1 2 3 4
// i=1  x x x x
//   2  x x x x
//   3  x x x x
//   4  x x x x
//  outer loop: j=3,
//  inner loop: i=1
//  array element accessed b[i-1][j] becomes b[0][3], which in turn is b[4][2]
//  due to linearized column-major storage of the 2-D array.
//  This causes loop-carried data dependence between j=2 and j=3.
//
//
//Data race pair: b[i][j]@67:13:W vs. b[i-1][j]@67:22:R

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var n, m, length, argCount int
	var args []string
	var b [][]float64
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

	n = length
	m = length

	b = make([][]float64, n)
	for i := range b {
		b[i] = make([]float64, m)
	}

	//$omp parallel do private(i)
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
			for j := start; j <= end; j++ {
				for i := 1; i <= m; i++ {
					b[i-1][j-1] = b[i-2][j-1]
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("b(50,50)= %f\n", b[49][49])
}
