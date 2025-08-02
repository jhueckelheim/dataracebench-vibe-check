//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Race condition due to anti-dependence within a loop offloaded to accelerators.
//Data race pair: a[i]@29:13:W vs. a[i+1]@29:20:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, length int
	var a []int

	length = 1000

	a = make([]int, length)

	for i = 1; i <= length; i++ {
		a[i-1] = i
	}

	//$omp target map(a)
	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := (length - 1) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= length-1; start += chunkSize {
		end := start + chunkSize - 1
		if end > length-1 {
			end = length - 1
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = a[i] + 1
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
	//$omp end target

	for i = 1; i <= length; i++ {
		fmt.Printf("Values for i and a(i) are: %d %d\n", i, a[i-1])
	}
}
