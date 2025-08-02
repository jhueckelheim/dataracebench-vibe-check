//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimensional array computation:
//Only one loop is associated with omp taskloop.
//The inner loop's loop iteration variable will be shared if it is shared in the enclosing context.
//Data race pairs (we allow multiple ones to preserve the pattern):
//  Write_set = {j@36:20 (implicit step +1)}
//  Read_set = {j@36:20, j@37:35}
//  Any pair from Write_set vs. Write_set  and Write_set vs. Read_set is a data race pair.

//need to run with large thread number and large num of iterations.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Module DRB095 translated to package-level variables
var a [][]int

func main() {
	var length int
	length = 100
	a = make([][]int, length)
	for i := range a {
		a[i] = make([]int, length)
	}

	//$omp parallel
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp single
		var wgSingle sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := length / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		//$omp taskloop
		for start := 1; start <= length; start += chunkSize {
			end := start + chunkSize - 1
			if end > length {
				end = length
			}
			wgSingle.Add(1)
			go func(start, end int) {
				defer wgSingle.Done()
				for i := start; i <= end; i++ {
					for j := 1; j <= length; j++ {
						a[i-1][j-1] = a[i-1][j-1] + 1
					}
				}
			}(start, end)
		}
		wgSingle.Wait()
		//$omp end taskloop
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	fmt.Printf("a(50,50) =%3d\n", a[49][49])
}
