//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// * Cover the implicitly determined rule: In an orphaned task generating construct,
// * formal arguments passed by reference are firstprivate.
// * This requires OpenMP 4.5 to work.
// * Earlier OpenMP does not allow a reference type for a variable within firstprivate().
// * No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable (module equivalent)
var a []int

func genTask(i int) {
	//$omp task
	// i is firstprivate (passed by reference in original)
	go func(iPrivate int) {
		a[iPrivate-1] = iPrivate + 1 // No race - i is firstprivate
	}(i)
	//$omp end task
}

func main() {
	a = make([]int, 100)

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			if threadID == 0 {
				var taskWg sync.WaitGroup
				for i := 1; i <= 100; i++ {
					taskWg.Add(1)
					go func(iCopy int) { // firstprivate equivalent
						defer taskWg.Done()
						a[iCopy-1] = iCopy + 1
					}(i)
				}
				taskWg.Wait()
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	for i := 1; i <= 100; i++ {
		if a[i-1] != i+1 {
			fmt.Printf("warning: a(%d) = %d not expected %d\n", i, a[i-1], i+1)
		}
	}
}
