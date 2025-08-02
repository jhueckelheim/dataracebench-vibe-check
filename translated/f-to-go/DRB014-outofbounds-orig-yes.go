//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The outmost loop is parallelized.
//But the inner level loop has out of bound access for b[i][j] when i equals to 1.
//This will case memory access of a previous columns's last element.
//
//For example, an array of 4x4:
//    j=1 2 3 4
// i=1  x x x x
//   2  x x x x
//   3  x x x x
//   4  x x x x
//  inner loop: i=1,
//  outer loop: j=3
//  array element accessed b[i-1][j] becomes b[0][3], which in turn is b[4][2]
//  due to linearized column-major storage of the 2-D array.
//  This causes loop-carried data dependence between j=2 and j=3.
//
//Data race pair: b[i][j]@41:13:W vs. b[i-1][j]@41:22:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, n, m int
	var b [][]float32

	n = 100
	m = 100

	b = make([][]float32, n)
	for i = 0; i < n; i++ {
		b[i] = make([]float32, m)
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
			var i int // private variable
			for j := start; j <= end; j++ {
				for i = 1; i <= m; i++ {
					// Out of bounds access when i=1: preserving the race condition from original
					b[i-1][j-1] = b[i-2][j-1] // This creates out-of-bounds access preserving the original race
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("b(50,50)= %v\n", b[49][49])

	// deallocate(b)
}
