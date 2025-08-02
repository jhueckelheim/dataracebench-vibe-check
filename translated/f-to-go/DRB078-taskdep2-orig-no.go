//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two tasks with depend clause to ensure execution order, no data races.
//i is shared for two tasks based on implicit data-sharing attribute rules.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var i int
	i = 0

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			//$omp single
			// Only one thread executes the single block
			if threadID == 0 {
				// Sequential execution of tasks with dependencies
				var taskWg sync.WaitGroup

				// First task (depend out:i)
				taskWg.Add(1)
				go func() {
					defer taskWg.Done()
					time.Sleep(3 * time.Second) // sleep(3)
					i = 3
				}()
				taskWg.Wait() // Wait for first task

				// Second task (depend out:i) - must execute after first
				taskWg.Add(1)
				go func() {
					defer taskWg.Done()
					i = 2
				}()
				taskWg.Wait() // Wait for second task
			}
			//$omp end single
		}(threadID)
	}
	wg.Wait()
	//$omp end parallel

	if i != 2 {
		fmt.Printf("%d\n", i)
	}
}
