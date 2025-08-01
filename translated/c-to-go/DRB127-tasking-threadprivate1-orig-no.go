/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* This example is referred from OpenMP Application Programming Interface 5.0, example tasking.7.c
 * A task switch may occur at a task scheduling point. A single thread may execute both of the
 * task regions that modify tp. The parts of these task regions in which tp is modified may be
 * executed in any order so the resulting value of var can be either 1 or 2.
 * There is a race pair but no data race due to threadprivate nature.
 */
package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Thread-local storage simulation using goroutine-local variables
var variable int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(1)

	// Outer task
	go func() {
		defer wg.Done()

		// Thread-private variable for this goroutine
		tp := 0

		var innerWg sync.WaitGroup
		innerWg.Add(1)

		// Inner task
		go func() {
			defer innerWg.Done()

			// This task modifies the thread-private tp
			tp = 1

			// Nested task (task scheduling point)
			var nestedWg sync.WaitGroup
			nestedWg.Add(1)
			go func() {
				defer nestedWg.Done()
				// Empty task - provides scheduling point
			}()
			nestedWg.Wait()

			// Read tp value after potential task switch
			variable = tp // tp is still 1 for this execution path
		}()

		// Concurrent modification in outer task
		tp = 2

		innerWg.Wait()
	}()

	wg.Wait()

	if variable == 2 {
		fmt.Printf("%d\n", variable)
	}
}
