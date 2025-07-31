/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* The assignment to a is not synchronized with the update of a as a result of the
 * reduction computation in the for loop.
 * Data Race pair: a (write vs. write)
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

	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Master thread initializes 'a'
			if threadID == 0 {
				atomic.StoreInt64(&a, 0) // Data race: write to a
			}

			// All threads participate in reduction - no barrier after master initialization
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

			// Data race: concurrent read-modify-write to 'a' while master may be writing
			atomic.AddInt64(&a, localSum) // Data race: reduction update conflicts with master write

			// Single thread prints result
			if threadID == 0 {
				fmt.Printf("Sum is %d\n", atomic.LoadInt64(&a))
			}
		}(t)
	}

	wg.Wait()
}
