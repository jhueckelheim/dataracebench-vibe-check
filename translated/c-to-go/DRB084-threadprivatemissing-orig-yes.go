/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB084-threadprivatemissing-orig-yes.c

Description: A file-scope variable used within a function called by a parallel region.
No threadprivate is used to avoid data races.

Original Data race pairs: sum0@61:3:W vs. sum0@61:8:R
                         sum0@61:3:W vs. sum0@61:3:W
*/

package main

import (
	"fmt"
	"sync"
)

var sum0 int = 0 // Global variable - shared across goroutines (causes data race)
var sum1 int = 0

func foo(i int) {
	// Data race: multiple goroutines read and write shared global sum0
	sum0 = sum0 + i
}

func main() {
	var sum int = 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Parallel region with work distribution
	numThreads := 4
	itemsPerThread := 1000 / numThreads

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			start := threadID*itemsPerThread + 1
			end := start + itemsPerThread
			if threadID == numThreads-1 {
				end = 1001 // Handle remainder for last thread
			}

			// Parallel for loop simulation
			for i := start; i < end; i++ {
				foo(i) // Data race occurs here
			}

			// Critical section for sum accumulation
			mu.Lock()
			sum = sum + sum0
			mu.Unlock()
		}(t)
	}

	wg.Wait()

	// Reference calculation
	for i := 1; i <= 1000; i++ {
		sum1 = sum1 + i
	}

	fmt.Printf("sum=%d; sum1=%d\n", sum, sum1)
	// Note: Due to data race, sum != sum1 is likely
}
