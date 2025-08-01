/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB103-master-orig-no.c

Description: A master directive is used to protect memory accesses.
Only the master thread (thread 0) executes the code block.
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var k int
	var wg sync.WaitGroup
	var once sync.Once

	numThreads := 5

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			// Simulate master directive - only first goroutine (master) executes
			// In OpenMP, master means thread 0, here we use sync.Once for similar effect
			once.Do(func() {
				k = numThreads // Simulate omp_get_num_threads()
				fmt.Printf("Number of Threads requested = %d\n", k)
			})
		}(t)
	}

	wg.Wait()
}
