/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB093-doall2-collapse-orig-no.c

Description: Two-dimensional array computation:
collapse(2) is used to associate two loops with omp for.
The corresponding loop iteration variables are private.
*/

package main

import (
	"sync"
)

var a [100][100]int

func main() {
	var wg sync.WaitGroup

	// Simulate OpenMP collapse(2) by flattening nested loops into single parallel loop
	totalIterations := 100 * 100
	numThreads := 4
	iterationsPerThread := totalIterations / numThreads

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			start := threadID * iterationsPerThread
			end := start + iterationsPerThread
			if threadID == numThreads-1 {
				end = totalIterations
			}

			// Each goroutine processes a range of flattened iterations
			// Both i and j are private to each goroutine (simulates collapse behavior)
			for iteration := start; iteration < end; iteration++ {
				i := iteration / 100 // Private i variable
				j := iteration % 100 // Private j variable
				a[i][j] = a[i][j] + 1
			}
		}(t)
	}

	wg.Wait()
}
