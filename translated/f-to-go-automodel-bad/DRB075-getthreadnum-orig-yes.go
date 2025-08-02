//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Test if the semantics of omp_get_thread_num() is correctly recognized.
//Thread with id 0 writes numThreads while other threads read it, causing data races.
//Data race pair: numThreads@22:9:W vs. numThreads@24:31:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var numThreads int
	numThreads = 0

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if true { // In Go, we can't easily get thread ID, so we simulate the condition
			numThreads = runtime.NumCPU()
		} else {
			fmt.Printf("numThreads = %d\n", numThreads)
		}
	}()
	wg.Wait()
	//$omp endparallel
}
