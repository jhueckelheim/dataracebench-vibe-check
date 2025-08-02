//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//tasks with depend clauses to ensure execution order, no data races.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var i, j, k int
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
				var taskWg sync.WaitGroup

				// Producer task (depend out:i)
				taskWg.Add(1)
				go func() {
					defer taskWg.Done()
					time.Sleep(3 * time.Second) // sleep(3)
					i = 1
				}()
				taskWg.Wait() // Wait for producer

				// Consumer tasks (depend in:i) - can execute in parallel after producer
				var consumerWg sync.WaitGroup

				consumerWg.Add(1)
				go func() {
					defer consumerWg.Done()
					j = i // No race - dependency ensures i is ready
				}()

				consumerWg.Add(1)
				go func() {
					defer consumerWg.Done()
					k = i // No race - dependency ensures i is ready
				}()

				consumerWg.Wait() // Wait for both consumers
			}
			//$omp end single
		}(threadID)
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("j = %3d  k = %3d\n", j, k)

	if j != 1 && k != 1 {
		fmt.Printf("Race Condition\n")
	}
}
