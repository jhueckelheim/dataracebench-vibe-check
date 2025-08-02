//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//omp_get_thread_num() is used to ensure serial semantics. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var numThreads int

	//$omp parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		if true { // In Go, we can't easily get thread ID, so we simulate the condition
			numThreads = runtime.NumCPU()
		}
		mu.Unlock()
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("numThreads =%3d\n", numThreads)
}
