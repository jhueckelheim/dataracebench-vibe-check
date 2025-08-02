//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//For the case of a variable which is referenced within a construct:
//static data member should be shared, unless it is within a threadprivate directive.
//
//Dependence pair: counter@37:5:W vs. counter@37:5:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variables (module equivalent)
var counter int    // Static shared variable
var pcounter int   // Would be threadprivate in original

type A struct {
	counter  int
	pcounter int
}

func main() {
	c := A{counter: 0, pcounter: 0}
	_ = c // Use c to avoid unused variable

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter = counter + 1   // RACE: Shared static variable modified without sync
			pcounter = pcounter + 1 // This would be threadprivate in original
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d %d\n", counter, pcounter)
}