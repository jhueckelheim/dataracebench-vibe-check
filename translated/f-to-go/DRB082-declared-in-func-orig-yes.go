//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A variable is declared inside a function called within a parallel region.
//The variable should be shared if it uses static storage.
//
//Data race pair: i@19:7:W vs. i@19:7:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable simulates static/save storage
var globalI int

func foo() {
	globalI = globalI + 1 // RACE: All goroutines access same static variable
	fmt.Printf("%d\n", globalI)
}

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			foo() // RACE: Static variable shared across all function calls
		}()
	}
	wg.Wait()
	//$omp end parallel
}
