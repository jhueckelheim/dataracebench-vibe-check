/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB083-declared-in-func-orig-no.c

Description: A variable is declared inside a function called within a parallel region.
The variable should be private if it does not use static storage.
*/

package main

import (
	"sync"
)

func foo() {
	// Local variable - each function call gets its own copy
	var q int = 0
	q += 1 // No data race: q is local to each function call
}

func main() {
	var wg sync.WaitGroup

	// Parallel region - multiple goroutines call foo()
	numThreads := 5
	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// No data race: each call to foo() has its own local variable q
			foo()
		}()
	}

	wg.Wait()
}
