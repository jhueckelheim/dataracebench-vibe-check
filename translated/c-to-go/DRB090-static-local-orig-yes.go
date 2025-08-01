/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB090-static-local-orig-yes.c

Description: For a variable declared in a scope inside an OpenMP construct:
* private if the variable has an automatic storage duration
* shared if the variable has a static storage duration.

Original Dependence pairs:
   tmp@73:7:W vs. tmp@73:7:W
   tmp@73:7:W vs. tmp@74:14:R
*/

package main

import (
	"fmt"
	"sync"
)

// Package-level variable to simulate C static variable
var staticTmp int

func main() {
	const len = 100
	a := make([]int, len)
	b := make([]int, len)

	// Initialize arrays
	for i := 0; i < len; i++ {
		a[i] = i
		b[i] = i
	}

	var wg1 sync.WaitGroup

	// First parallel region - static storage simulation (data race)
	numThreads := 4
	itemsPerThread := len / numThreads

	for t := 0; t < numThreads; t++ {
		wg1.Add(1)
		go func(threadID int) {
			defer wg1.Done()

			start := threadID * itemsPerThread
			end := start + itemsPerThread
			if threadID == numThreads-1 {
				end = len
			}

			// Data race: all goroutines access shared staticTmp variable
			for i := start; i < end; i++ {
				staticTmp = a[i] + i // Data race here
				a[i] = staticTmp     // and here
			}
		}(t)
	}
	wg1.Wait()

	var wg2 sync.WaitGroup

	// Second parallel region - automatic storage (no data race)
	for t := 0; t < numThreads; t++ {
		wg2.Add(1)
		go func(threadID int) {
			defer wg2.Done()

			start := threadID * itemsPerThread
			end := start + itemsPerThread
			if threadID == numThreads-1 {
				end = len
			}

			// No data race: each goroutine has its own local tmp variable
			var tmp int // Local variable - automatic storage
			for i := start; i < end; i++ {
				tmp = b[i] + i
				b[i] = tmp
			}
		}(t)
	}
	wg2.Wait()

	fmt.Printf("a[50]=%d b[50]=%d\n", a[50], b[50])
}
