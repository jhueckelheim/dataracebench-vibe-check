//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A function argument passed by value should be private inside the function.
//Variable i is read only. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func f1(i int) { // Pass by value - i is private copy in each function call
	i = i + 1 // No race - each goroutine has its own copy
}

func main() {
	var i int
	i = 0

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f1(i) // No race - pass by value makes i private within f1
		}()
	}
	wg.Wait()
	//$omp end parallel

	if i != 0 {
		fmt.Printf("i = %3d\n", i)
	}
}
