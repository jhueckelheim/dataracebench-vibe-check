//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// argument pass-by-reference
// its data-sharing attribute is the same as its actual argument's. i and q are shared.
// Data race pair: q@15:5:W vs. q@15:5:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func f1(q *int) {
	*q = *q + 1 // RACE: Multiple goroutines modifying shared variable through pointer
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
			f1(&i) // RACE: Pass by reference - all threads modify same variable
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("i = %d\n", i)
}
