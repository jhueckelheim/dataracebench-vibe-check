/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB075-getthreadnum-orig-yes.c

Description: Test if the semantics of omp_get_thread_num() is correctly recognized.
Thread with id 0 writes numThreads while other threads read it, causing data races.

Original Data race pair: numThreads@60:7:W vs. numThreads@64:33:R
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var numThreads int = 0
	var wg sync.WaitGroup

	// Simulate parallel region with multiple goroutines
	for threadID := 0; threadID < 5; threadID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if id == 0 {
				// Thread 0 writes to numThreads - this creates the data race
				numThreads = 5 // Total number of threads
			} else {
				// Other threads read from numThreads - data race with the write above
				fmt.Printf("Thread %d sees numThreads=%d\n", id, numThreads)
			}
		}(threadID)
	}

	wg.Wait()
	fmt.Printf("Final numThreads=%d\n", numThreads)
}
