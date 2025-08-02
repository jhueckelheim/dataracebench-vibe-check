//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This benchmark is extracted from flush_nolist.1c of OpenMP Application
//Programming Interface Examples Version 4.5.0 .
//We added one critical section to make it a test with only one pair of data races.
//The data race will not generate wrong result though. So the assertion always passes.
//Data race pair:  i@37:13:W vs. i@38:15:R

package main

import (
	"fmt"
	"sync"
)

// Module DRB074 translated to package-level functions
func f1(q *int) {
	var mu sync.Mutex
	//$omp critical
	mu.Lock()
	*q = 1
	mu.Unlock()
	//$omp end critical
	//$omp flush
}

func main() {
	var i, sum int
	i = 0
	sum = 0

	//$omp parallel reduction(+:sum) num_threads(10)
	var wg sync.WaitGroup
	var sumMutex sync.Mutex
	numThreads := 10
	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		go func() {
			defer wg.Done()
			f1(&i)
			localSum := i
			sumMutex.Lock()
			sum += localSum
			sumMutex.Unlock()
		}()
	}
	wg.Wait()
	//$omp end parallel

	if sum != 10 {
		fmt.Printf("sum = %d\n", sum)
	}
}
