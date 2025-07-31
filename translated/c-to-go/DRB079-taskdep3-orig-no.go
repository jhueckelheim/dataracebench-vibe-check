/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB079-taskdep3-orig-no.c

Description: Tasks with depend clauses to ensure execution order, no data races.
One task produces a value, two tasks consume it.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// Simulated delay function
func delay(microseconds int) {
	time.Sleep(time.Duration(microseconds) * time.Microsecond)
}

func main() {
	var i, j, k int = 0, 0, 0
	var wg sync.WaitGroup

	// Channel to coordinate dependencies
	taskComplete := make(chan bool, 2)

	// Producer task - writes i=1 with delay
	wg.Add(1)
	go func() {
		defer wg.Done()
		delay(10000) // 10ms delay
		i = 1
		// Signal completion to both consumer tasks
		taskComplete <- true
		taskComplete <- true
	}()

	// First consumer task - reads i into j
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-taskComplete // Wait for producer to complete
		j = i
	}()

	// Second consumer task - reads i into k
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-taskComplete // Wait for producer to complete
		k = i
	}()

	wg.Wait()

	fmt.Printf("j=%d k=%d\n", j, k)

	// Assertion
	if j != 1 || k != 1 {
		panic(fmt.Sprintf("Assertion failed: expected j=1 and k=1, got j=%d k=%d", j, k))
	}
}
