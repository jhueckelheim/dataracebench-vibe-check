/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB096-doall2-taskloop-collapse-orig-no.c

Description: Two-dimensional array computation:
Two loops are associated with omp taskloop due to collapse(2).
Both loop index variables are private.
taskloop requires OpenMP 4.5 compilers.
*/

package main

import (
	"fmt"
	"sync"
)

var a [100][100]int

func main() {
	var wg sync.WaitGroup

	// Simulate OpenMP taskloop with collapse(2)
	// Both i and j are private due to collapse flattening
	totalIterations := 100 * 100

	// Create tasks for chunks of iterations
	numTasks := 10
	iterationsPerTask := totalIterations / numTasks

	for task := 0; task < numTasks; task++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			start := taskID * iterationsPerTask
			end := start + iterationsPerTask
			if taskID == numTasks-1 {
				end = totalIterations
			}

			// Each task processes a range of flattened iterations
			// Both i and j are private to each task (simulates collapse behavior)
			for iteration := start; iteration < end; iteration++ {
				i := iteration / 100 // Private i variable
				j := iteration % 100 // Private j variable
				a[i][j] += 1
			}
		}(task)
	}

	wg.Wait()

	fmt.Printf("a[50][50]=%d\n", a[50][50])
}
