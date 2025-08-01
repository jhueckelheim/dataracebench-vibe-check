/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB080-func-arg-orig-yes.c

Description: A function argument is passed by reference:
its data-sharing attribute is the same as its actual argument's.
i is shared. *q is shared.

Original Data race pair: *q@59:4:W vs. *q@59:4:W
*/

package main

import (
	"fmt"
	"sync"
)

// Function that increments the value pointed to by q
func f1(q *int) {
	*q += 1 // Data race: multiple goroutines writing to same memory location
}

func main() {
	var i int = 0
	var wg sync.WaitGroup

	// Parallel region - multiple goroutines call f1 with same shared variable
	numThreads := 5
	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Data race: all goroutines call f1 with pointer to same variable i
			f1(&i)
		}()
	}

	wg.Wait()

	fmt.Printf("i=%d\n", i)

	// Note: Due to data race, the final value of i is unpredictable
	// It could be anywhere from 1 to numThreads depending on the interleaving
}
