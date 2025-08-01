/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB074-flush-orig-yes.c

Description: This benchmark is extracted from flush_nolist.1c of OpenMP Application
Programming Interface Examples Version 4.5.0.
We added one critical section to make it a test with only one pair of data races.
The data race will not generate wrong result though. So the assertion always passes.

Original Data race pair: *q@60:3:W vs. i@71:11:R
*/

package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

func f1(q *int) {
	// Critical section - only one goroutine can execute this at a time
	mu.Lock()
	*q = 1
	mu.Unlock()
	// Note: Go's memory model provides stronger guarantees than C/OpenMP flush
}

func main() {
	var i int = 0
	var sum int = 0
	var wg sync.WaitGroup
	var sumMu sync.Mutex

	// Parallel execution with 10 goroutines
	for t := 0; t < 10; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Call f1 which sets i=1 under critical section
			f1(&i)

			// Data race: reading i here while another goroutine might be writing to it
			// Even though f1 uses a critical section, the read of i here is not protected
			localI := i

			// Protected sum update
			sumMu.Lock()
			sum += localI
			sumMu.Unlock()
		}()
	}

	wg.Wait()

	// Assertion should always pass despite the data race
	if sum != 10 {
		panic(fmt.Sprintf("Assertion failed: expected sum=10, got sum=%d", sum))
	}
	fmt.Printf("sum=%d\n", sum)
}
