//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A master directive is used to protect memory accesses. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var k int

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//$omp master
		k = runtime.NumCPU()
		fmt.Printf("Number of threads requested =%8d\n", k)
		//$omp end master
	}()
	wg.Wait()
	//$omp end parallel
}
