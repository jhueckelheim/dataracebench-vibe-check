//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This loop has loop-carried output-dependence due to x=... at line 21.
//The problem can be solved by using lastprivate(x).
//Data race pair: x@21:9:W vs. x@21:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var x, length int
	length = 10000

	//$omp parallel do private(i)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := (length + 1) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 0; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			var i int // private variable
			for i = start; i <= end; i++ {
				x = i // This creates the data race - x is shared, not lastprivate
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("x = %v\n", x)
}
