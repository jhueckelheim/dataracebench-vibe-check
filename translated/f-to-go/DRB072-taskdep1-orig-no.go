//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two tasks with depend clause to ensure execution order:
//i is shared for two tasks based on implicit data-sharing attribute rules. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

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

			//$omp single
			// Simulate single construct - only one goroutine executes this
			if threadID == 0 {
				// Task dependency simulation: first task (depend out:i)
				var taskWg sync.WaitGroup

				taskWg.Add(1)
				go func() {
					defer taskWg.Done()
					i = 1 // First task sets i
				}()
				taskWg.Wait() // Wait for first task to complete

				// Second task (depend in:i) - executes after first task
				taskWg.Add(1)
				go func() {
					defer taskWg.Done()
					i = 2 // Second task reads i (dependency satisfied) and sets it
				}()
				taskWg.Wait() // Wait for second task to complete
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	if i != 2 {
		fmt.Printf("i is not equal to 2\n")
	}
}
