/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB077-single-orig-no.c

Description: A single directive is used to protect a write.
The OpenMP single directive ensures only one thread executes the code block.
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int = 0
	var wg sync.WaitGroup
	var once sync.Once

	// Parallel region with multiple goroutines
	numThreads := 5
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Use sync.Once to simulate OpenMP single directive
			// Only one goroutine will execute this function
			once.Do(func() {
				count += 1
			})
		}()
	}

	wg.Wait()

	fmt.Printf("count= %d\n", count)
}
