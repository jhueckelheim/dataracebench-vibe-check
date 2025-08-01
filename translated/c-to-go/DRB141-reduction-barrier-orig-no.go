/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* To avoid data race, the initialization of the original list item "a" should complete before any
 * update of a as a result of the reduction clause. This can be achieved by adding an explicit
 * barrier after the assignment a=0, or by enclosing the assignment a=0 in a single directive
 * or by initializing a before the start of the parallel region.
 * */
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var a int64
	numThreads := runtime.NumCPU()
	var wg sync.WaitGroup
	var barrier sync.WaitGroup

	// Set up barrier for all threads
	barrier.Add(numThreads)

	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Master thread initializes 'a'
			if threadID == 0 {
				atomic.StoreInt64(&a, 0)
			}

			// Explicit barrier - wait for initialization to complete
			barrier.Done()
			barrier.Wait()

			// All threads participate in reduction after barrier
			localSum := int64(0)
			chunkSize := 10 / numThreads
			start := threadID * chunkSize
			end := start + chunkSize
			if threadID == numThreads-1 {
				end = 10
			}

			for i := start; i < end; i++ {
				localSum += int64(i)
			}

			// Safe: reduction happens after initialization is complete
			atomic.AddInt64(&a, localSum)
		}(t)
	}

	wg.Wait()

	// Print result after all threads complete
	fmt.Printf("Sum is %d\n", atomic.LoadInt64(&a))
}
