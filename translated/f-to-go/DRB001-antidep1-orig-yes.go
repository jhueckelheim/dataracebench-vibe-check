//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A loop with loop-carried anti-dependence.
//Data race pair: a[i+1]@25:9:W vs. a[i]@25:16:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, len int
	var a [1000]int

	len = 1000

	for i = 1; i <= len; i++ {
		a[i-1] = i
	}

	//$omp parallel do
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
			for i := start; i <= end; i++ {
				a[i-1] = a[i] + 1
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("a(500)=%3d\n", a[499])
}
