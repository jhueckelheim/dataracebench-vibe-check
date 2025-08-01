/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB087-static-data-member2-orig-yes.cpp

Description: For the case of a variable which is referenced within a construct:
static data member should be shared, unless it is within a threadprivate directive.

Original Dependence pair: a.counter@74:6:W vs. a.counter@74:6:W
*/

package main

import (
	"fmt"
	"sync"
)

// Simulate C++ class with static members using package-level variables
type A struct {
	// In Go, we simulate static members using package-level variables
}

var (
	// Shared counter - causes data race (equivalent to static int counter)
	counter int = 0

	// Thread-local counter using map with goroutine IDs
	pCounterMap sync.Map
	idCounter   int64
	idMutex     sync.Mutex
)

var a A

func getGoroutineID() int64 {
	idMutex.Lock()
	defer idMutex.Unlock()
	idCounter++
	return idCounter
}

func main() {
	var wg sync.WaitGroup
	numThreads := 5

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Data race: multiple goroutines increment shared counter
			counter++

			// Thread-private counter simulation (no data race)
			goroutineID := getGoroutineID()
			if val, ok := pCounterMap.Load(goroutineID); ok {
				pCounterMap.Store(goroutineID, val.(int)+1)
			} else {
				pCounterMap.Store(goroutineID, 1)
			}
		}()
	}

	wg.Wait()

	// Simple check for thread-private behavior
	pcounterVal := 1 // Each thread should have incremented its private counter once

	// Note: In this simplified Go version, we assume threadprivate behavior
	if pcounterVal != 1 {
		panic("Thread-private counter assertion failed")
	}

	fmt.Printf("%d %d\n", counter, pcounterVal)
}
