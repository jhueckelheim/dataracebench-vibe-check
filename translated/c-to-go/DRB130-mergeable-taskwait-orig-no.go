/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * Taken from OpenMP Examples 5.0, example tasking.12.c
 * x is a shared variable the outcome does not depend on whether or not the task is merged (that is,
 * the task will always increment the same variable and will always compute the same value for x).
 */
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var x int64 = 2
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(1)

	// Simulate mergeable task behavior with shared variable
	// Regardless of execution context, x is shared and behavior is consistent
	if rand.Intn(2) == 0 {
		// "Merged" execution - task executes in calling goroutine context
		go func() {
			defer wg.Done()
			atomic.AddInt64(&x, 1) // Atomic access to shared variable
		}()
	} else {
		// "Not merged" execution - task executes in separate goroutine
		go func() {
			defer wg.Done()
			atomic.AddInt64(&x, 1) // Same atomic access to shared variable
		}()
	}

	wg.Wait()

	fmt.Printf("%d\n", atomic.LoadInt64(&x))
}
