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

var globalI int // Global variable to be shared across goroutines

func f1(iPtr *int) {
	var criticalMutex sync.Mutex

	//$omp critical
	criticalMutex.Lock()
	*iPtr = 1
	criticalMutex.Unlock()
	//$omp end critical

	//$omp flush
	// Go's memory model handles flush semantics automatically
}

func main() {
	var sum int
	globalI = 0
	sum = 0

	//$omp parallel reduction(+:sum) num_threads(10)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numThreads := 10

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localSum := 0

			f1(&globalI)       // Write to globalI (through critical section)
			localSum = globalI // RACE: Read globalI without synchronization!

			// Reduction
			mu.Lock()
			sum += localSum
			mu.Unlock()
		}()
	}
	wg.Wait()
	//$omp end parallel

	if sum != 10 {
		fmt.Printf("sum = %d\n", sum)
	}
}
