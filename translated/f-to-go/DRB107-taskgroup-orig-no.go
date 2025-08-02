//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//* This is a program based on a test contributed by Yizi Gu@Rice Univ.
//* Use taskgroup to synchronize two tasks. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var result int
	result = 0

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			if threadID == 0 { // Only one thread executes single
				//$omp taskgroup
				var taskGroupWg sync.WaitGroup

				//$omp task
				taskGroupWg.Add(1)
				go func() {
					defer taskGroupWg.Done()
					time.Sleep(3 * time.Second) // sleep(3)
					result = 1
				}()

				taskGroupWg.Wait() // taskgroup ensures this task completes
				//$omp end taskgroup

				//$omp task
				// This task runs after taskgroup completes
				go func() {
					result = 2 // No race - runs after first task due to taskgroup
				}()
				//$omp end task
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("result = %8d\n", result)
}
