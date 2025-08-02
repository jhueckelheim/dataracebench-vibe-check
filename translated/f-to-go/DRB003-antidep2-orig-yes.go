//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A two-level loop nest with loop carried anti-dependence on the outer level.
//Data race pair: a[i][j]@29:13:W vs. a[i+1][j]@29:31:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, j, len int
	var a [20][20]float32

	len = 20

	for i = 1; i <= len; i++ {
		for j = 1; j <= len; j++ {
			a[i-1][j-1] = 0.5
		}
	}

	//$omp parallel do private(j)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := (len - 1) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= len-1; start += chunkSize {
		end := start + chunkSize - 1
		if end > len-1 {
			end = len - 1
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			var j int // private variable
			for i := start; i <= end; i++ {
				for j = 1; j <= len; j++ {
					a[i-1][j-1] = a[i-1][j-1] + a[i][j-1]
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("a(10,10) = %v\n", a[9][9])
}
