/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB078-taskdep2-orig-no.c

Description: Two tasks with depend clause to ensure execution order, no data races.
i is shared for two tasks based on implicit data-sharing attribute rules.
*/

package main

import (
	"fmt"
	"time"
)

// Simulated delay function
func delay(microseconds int) {
	time.Sleep(time.Duration(microseconds) * time.Microsecond)
}

func main() {
	var i int = 0

	// Channel to synchronize task execution order
	taskComplete := make(chan bool, 1)

	// First task - writes i=1 with delay
	go func() {
		delay(10000) // 10ms delay
		i = 1
		taskComplete <- true // Signal completion
	}()

	// Second task - waits for first task completion, then writes i=2
	go func() {
		<-taskComplete // Wait for first task to complete
		i = 2
	}()

	// Wait a bit for tasks to complete
	time.Sleep(20 * time.Millisecond)

	// Assertion
	if i != 2 {
		panic(fmt.Sprintf("Assertion failed: expected i=2, got i=%d", i))
	}

	fmt.Printf("Final i=%d\n", i)
}
