/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB081-func-arg-orig-no.c

Description: A function argument passed by value should be private inside the function.
Variable i is read only.
*/

package main

import (
	"fmt"
	"sync"
)

// Function that receives argument by value (copy) - no data race
func f1(q int) {
	q += 1 // Modifying the local copy, not the original
}

func main() {
	var i int = 0
	var wg sync.WaitGroup

	// Parallel region - multiple goroutines call f1 with value copy
	numThreads := 5
	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Pass by value - each goroutine gets its own copy of i
			f1(i) // No data race: i is only read, q is local copy
		}()
	}

	wg.Wait()

	// Assertion should pass - i is unchanged
	if i != 0 {
		panic(fmt.Sprintf("Assertion failed: expected i=0, got i=%d", i))
	}
	fmt.Printf("i=%d\n", i)
}
