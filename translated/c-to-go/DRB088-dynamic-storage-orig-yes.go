/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB088-dynamic-storage-orig-yes.c

Description: For the case of a variable which is not referenced within a construct:
objects with dynamic storage duration should be shared.
Putting it within a threadprivate directive may cause seg fault since
threadprivate copies are not allocated!

Original Dependence pair: *counter@63:6:W vs. *counter@63:6:W
*/

package main

import (
	"fmt"
	"sync"
)

var counter *int

func foo() {
	// Data race: multiple goroutines increment same shared memory location
	*counter++
}

func main() {
	// Dynamically allocate memory (equivalent to malloc)
	counter = new(int) // Go's new() is equivalent to malloc + initialization
	*counter = 0

	var wg sync.WaitGroup
	numThreads := 5

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Data race: all goroutines access same dynamically allocated memory
			foo()
		}()
	}

	wg.Wait()

	fmt.Printf("%d\n", *counter)

	// Note: Due to data race, final value is unpredictable
	// In Go, we don't need explicit free() - garbage collector handles it
}
