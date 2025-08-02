//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The thread encountering the taskwait directive at line 22 only waits for its child task
//(line 14-21) to complete. It does not wait for its descendant tasks (line 16-19). Data Race pairs, sum@36:13:W vs. sum@36:13:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var a, psum []int
	var sum int

	a = make([]int, 4)
	psum = make([]int, 4)

	//$omp parallel num_threads(2)
	var wg sync.WaitGroup
	numThreads := 2

	for threadID := 0; threadID < numThreads; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp do schedule(dynamic, 1)
			// Simple work distribution
			for i := 1; i <= 4; i++ {
				a[i-1] = i
			}
			//$omp end do

			//$omp single
			if threadID == 0 { // Only one thread executes single
				//$omp task
				var childWg sync.WaitGroup
				childWg.Add(1)
				go func() {
					defer childWg.Done()

					//$omp task (descendant task)
					go func() {
						psum[1] = a[2] + a[3] // This runs independently
					}()

					psum[0] = a[0] + a[1] // Child task work
				}()

				//$omp taskwait
				childWg.Wait() // Only waits for child, NOT descendant!

				sum = psum[1] + psum[0] // RACE: psum[1] may not be ready yet
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("sum = %d\n", sum)

	// deallocate(a,psum)
}
