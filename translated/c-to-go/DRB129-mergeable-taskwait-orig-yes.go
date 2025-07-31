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
 * The created task will access different instances of the variable x if the task is not merged,
 * as x is firstprivate, but it will access the same variable x if the task is merged. It can
 * print two different values for x depending on the decisions taken by the implementation.
 * Data Race Pairs: x (write vs. write) - behavior depends on goroutine execution context
 */
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	x := 2
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(1)

	// Simulate mergeable task behavior - task may execute in different contexts
	if rand.Intn(2) == 0 {
		// "Merged" execution - task executes in calling goroutine context
		go func() {
			defer wg.Done()
			// Direct access to parent's x variable
			x++ // Data race: direct modification of parent variable
		}()
	} else {
		// "Not merged" execution - task gets its own copy (firstprivate)
		go func() {
			defer wg.Done()
			// Task gets its own copy of x (firstprivate behavior)
			localX := x // Copy of x
			localX++
			// This doesn't affect the original x - different behavior!
			x = localX // Data race: concurrent write to shared x
		}()
	}

	wg.Wait()

	fmt.Printf("%d\n", x)
}
