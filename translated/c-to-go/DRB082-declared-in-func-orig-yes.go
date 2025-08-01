/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB082-declared-in-func-orig-yes.c

Description: A variable is declared inside a function called within a parallel region.
The variable should be shared if it uses static storage.

Original Data race pair: q@57:3:W vs. q@57:3:W
*/

package main

import (
	"sync"
)

// Package-level variable to simulate C static variable inside function
var q int

func foo() {
	// Data race: multiple goroutines access shared package variable q
	q += 1
}

func main() {
	var wg sync.WaitGroup

	// Parallel region - multiple goroutines call foo()
	numThreads := 5
	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Data race: all goroutines modify same shared variable q
			foo()
		}()
	}

	wg.Wait()

	// Note: Due to data race, final value of q is unpredictable
	// (could be anywhere from 1 to numThreads)
}
