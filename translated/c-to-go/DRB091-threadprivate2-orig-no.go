/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB091-threadprivate2-orig-no.c

Description: A file-scope variable used within a function called by a parallel region.
Use threadprivate to avoid data races.
This is the case for a variable referenced within a construct.
*/

package main

import (
	"fmt"
	"sync"
)

var sum1 int = 0

func main() {
	const len = 1000
	var sum int = 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Parallel region with thread-local sum0
	numThreads := 4
	itemsPerThread := len / numThreads

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			// Each goroutine has its own private sum0 (simulates threadprivate)
			var sum0 int = 0 // copyin(sum0) behavior - each thread starts with 0

			start := threadID * itemsPerThread
			end := start + itemsPerThread
			if threadID == numThreads-1 {
				end = len
			}

			// Parallel for loop simulation with direct access
			for i := start; i < end; i++ {
				sum0 = sum0 + i // No data race: each thread has private sum0
			}

			// Critical section for sum accumulation
			mu.Lock()
			sum = sum + sum0
			mu.Unlock()
		}(t)
	}

	wg.Wait()

	// Reference calculation
	for i := 0; i < len; i++ {
		sum1 = sum1 + i
	}

	fmt.Printf("sum=%d; sum1=%d\n", sum, sum1)

	// Assertion should pass
	if sum != sum1 {
		panic(fmt.Sprintf("Assertion failed: expected sum=%d, got sum=%d", sum1, sum))
	}
}
