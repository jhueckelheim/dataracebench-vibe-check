/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * Referred from worksharing_critical.1.c
 * A single thread executes the one and only section in the sections region, and executes the
 * critical region. The same thread encounters the nested parallel region, creates a new team
 * of threads, and becomes the master of the new team. One of the threads in the new team enters
 * the single region and increments i by 1. At the end of this example i is equal to 2.
 */
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	i := 1
	var outerWg sync.WaitGroup
	var criticalMutex sync.Mutex

	outerWg.Add(1)

	// Parallel sections region
	go func() {
		defer outerWg.Done()

		// Section 1 (only one section, so only one goroutine executes)
		criticalMutex.Lock() // Critical region
		{
			// Nested parallel region
			var innerWg sync.WaitGroup
			innerWg.Add(1)

			// Single execution within nested parallel region
			go func() {
				defer innerWg.Done()
				i++ // Only one goroutine increments i
			}()

			innerWg.Wait()
		}
		criticalMutex.Unlock()
	}()

	outerWg.Wait()

	fmt.Printf("%d\n", i) // Should print 2
}
