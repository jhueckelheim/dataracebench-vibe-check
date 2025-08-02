//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//For the case of a variable which is not referenced within a construct:
//objects with dynamic storage duration should be shared.
//Putting it within a threadprivate directive may cause seg fault since
// threadprivate copies are not allocated!
//
//Dependence pair: *counter@22:9:W vs. *counter@22:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level pointer (dynamic storage)
var counter *int

func foo() {
	*counter = *counter + 1 // RACE: Multiple threads modifying same dynamic storage
}

func main() {
	// Allocate dynamic storage
	counter = new(int)
	*counter = 0

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			foo() // RACE: All threads call foo which modifies shared *counter
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d\n", *counter)

	// deallocate(counter) - handled by Go's garbage collector
}
