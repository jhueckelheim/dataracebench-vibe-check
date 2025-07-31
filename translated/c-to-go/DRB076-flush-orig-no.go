/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB076-flush-orig-no.c

Description: This benchmark is extracted from flush_nolist.1c of OpenMP
Application Programming Interface Examples Version 4.5.0.

We privatize variable i to fix data races in the original example.
Once i is privatized, flush is no longer needed.
*/

package main

import (
	"fmt"
	"sync"
)

func f1(q *int) {
	*q = 1
}

func main() {
	var sum int = 0
	var wg sync.WaitGroup
	var sumMu sync.Mutex

	// Parallel execution with 10 goroutines
	for t := 0; t < 10; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Each goroutine has its own private variable i
			var i int = 0

			// Call f1 with private variable - no data race
			f1(&i)

			// Protected sum update using local variable
			sumMu.Lock()
			sum += i
			sumMu.Unlock()
		}()
	}

	wg.Wait()

	// Assertion should always pass
	if sum != 10 {
		panic(fmt.Sprintf("Assertion failed: expected sum=10, got sum=%d", sum))
	}
	fmt.Printf("sum=%d\n", sum)
}
