/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB099-targetparallelfor2-orig-no.c

Description: Use of omp target + map + array sections derived from pointers
Target offloading with memory mapping and parallel execution.
*/

package main

import (
	"fmt"
	"sync"
)

func foo(a, b []float64, N int) {
	// Simulate target offloading with memory mapping
	// In Go, slices are already reference types (similar to mapped arrays)

	var wg sync.WaitGroup
	numThreads := 4
	itemsPerThread := N / numThreads

	// Parallel for loop simulation
	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			start := threadID * itemsPerThread
			end := start + itemsPerThread
			if threadID == numThreads-1 {
				end = N
			}

			// Each thread processes its assigned range
			for i := start; i < end; i++ {
				b[i] = a[i] * float64(i)
			}
		}(t)
	}

	wg.Wait()
}

func main() {
	const len = 1000
	a := make([]float64, len)
	b := make([]float64, len)

	// Initialize arrays
	for i := 0; i < len; i++ {
		a[i] = float64(i) / 2.0
		b[i] = 0.0
	}

	// Call function with slice arguments (automatic "mapping")
	foo(a, b, len)

	fmt.Printf("b[50]=%f\n", b[50])
}
