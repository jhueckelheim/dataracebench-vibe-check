/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB104-nowait-barrier-orig-no.c

Description: This example is based on one code snippet extracted from a paper:
Ma etc. Symbolic Analysis of Concurrency Errors in OpenMP Programs, ICPP 2013

Explicit barrier to counteract nowait.
The nowait clause removes the implicit barrier, but explicit barrier ensures synchronization.
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	const len = 1000
	var error int
	a := make([]int, len)
	b := 5

	// Initialize array
	for i := 0; i < len; i++ {
		a[i] = i
	}

	var wg sync.WaitGroup
	var once sync.Once
	numThreads := 4
	itemsPerThread := len / numThreads

	// Channel to coordinate barrier
	barrierChan := make(chan bool, numThreads)

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			// Parallel for with nowait simulation (no implicit barrier)
			start := threadID * itemsPerThread
			end := start + itemsPerThread
			if threadID == numThreads-1 {
				end = len
			}

			for i := start; i < end; i++ {
				a[i] = b + a[i]*5
			}

			// Explicit barrier - signal completion
			barrierChan <- true

			// Wait for all threads to complete
			for i := 0; i < numThreads; i++ {
				<-barrierChan
			}

			// Single thread executes this (simulating omp single)
			once.Do(func() {
				error = a[9] + 1
			})
		}(t)
	}

	wg.Wait()

	// Assertion
	if error != 51 {
		panic(fmt.Sprintf("Assertion failed: expected error=51, got error=%d", error))
	}

	fmt.Printf("error = %d\n", error)
}
