//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//One dimension array computation
//with finer granularity than traditional 4 bytes.
//There is no data race pair.

package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	var a []string
	
	a = make([]string, 100)

	//$omp parallel do private(str)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := 100 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= 100; start += chunkSize {
		end := start + chunkSize - 1
		if end > 100 {
			end = 100
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				str := strconv.Itoa(i) // str is private to each goroutine
				a[i-1] = str           // No race - each thread works on different elements
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("a(i) %s\n", a[22])
}