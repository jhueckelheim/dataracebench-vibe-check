/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* The explicit flush directive that provides release semantics is needed
 * here to complete the synchronization. No data race.
 * */
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	runtime.GOMAXPROCS(2) // Force exactly 2 threads

	x := 0
	var y int64 = 0
	var criticalMutex sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	go func() { // Thread 0
		defer wg.Done()

		criticalMutex.Lock()
		x = 10
		criticalMutex.Unlock()

		// Explicit memory barrier (flush) - ensures x write is visible
		runtime.Gosched() // Force context switch to ensure visibility

		atomic.StoreInt64(&y, 1) // Atomic write with release semantics
	}()

	go func() { // Thread 1
		defer wg.Done()

		var tmp int64 = 0
		for tmp == 0 {
			tmp = atomic.LoadInt64(&y) // Atomic read with acquire semantics
		}

		criticalMutex.Lock()
		if x != 10 { // Safe: x write is guaranteed to be visible after acquire
			fmt.Printf("x = %d\n", x)
		}
		criticalMutex.Unlock()
	}()

	wg.Wait()
}
