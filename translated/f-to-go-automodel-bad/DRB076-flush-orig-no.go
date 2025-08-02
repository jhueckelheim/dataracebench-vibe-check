//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This benchmark is extracted from flush_nolist.1c of OpenMP
//Application Programming Interface Examples Version 4.5.0 .
//
//We privatize variable i to fix data races in the original example.
//Once i is privatized, flush is no longer needed. No data race pairs.

package main

import (
	"fmt"
	"sync"
)

// Module DRB076 translated to package-level functions
func f1(q *int) {
	*q = 1
}

func main() {
	var sum int
	sum = 0

	//$omp parallel reduction(+:sum) num_threads(10) private(i)
	var wg sync.WaitGroup
	var sumMutex sync.Mutex
	numThreads := 10
	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		go func() {
			defer wg.Done()
			var i int = 0
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
