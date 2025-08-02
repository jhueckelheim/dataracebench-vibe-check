//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A linear expression is used as array subscription.
//Data race pair: a[2*i]@27:9:W vs. a[i]@27:18:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var length int
	var a []int

	length = 2000
	a = make([]int, length)

	for i := 1; i <= length; i++ {
		a[i-1] = i
	}

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 1000 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= 1000; start += chunkSize {
		end := start + chunkSize - 1
		if end > 1000 {
			end = 1000
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[2*i-1] = a[i-1] + 1 // Linear indexing race condition
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("a(1002) = %3d\n", a[1001])

	// deallocate(a)
}
