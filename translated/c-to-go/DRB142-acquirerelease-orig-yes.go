/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* The below program will fail to order the write to x on thread 0 before the read from x on thread 1.
 * The implicit release flush on exit from the critical region will not synchronize with the acquire
 * flush that occurs on the atomic read operation performed by thread 1. This is because implicit
 * release flushes that occur on a given construct may only synchronize with implicit acquire flushes
 * on a compatible construct (and vice-versa) that internally makes use of the same synchronization
 * variable.
 *
 * Missing memory barrier between critical section and atomic operation causes data race.
 * Data Race pair: x (write vs. read)
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
		x = 10 // Data race: write to x
		criticalMutex.Unlock()

		// Missing memory barrier here - no guarantee that x write is visible
		// to other thread before y is set

		atomic.StoreInt64(&y, 1) // Atomic write
	}()

	go func() { // Thread 1
		defer wg.Done()

		var tmp int64 = 0
		for tmp == 0 {
			tmp = atomic.LoadInt64(&y) // Atomic read with acquire semantics
		}

		criticalMutex.Lock()
		if x != 10 { // Data race: read x which may not be visible yet
			fmt.Printf("x = %d\n", x)
		}
		criticalMutex.Unlock()
	}()

	wg.Wait()
}
