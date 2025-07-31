/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * The scheduling constraints prohibit a thread in the team from executing
 * a new task that modifies tp while another such task region tied to
 * the same thread is suspended. Therefore, the value written will
 * persist across the task scheduling point.
 * No Data Race due to scheduling constraints
 */
package main

import (
	"runtime"
	"sync"
)

// Thread-local storage simulation
var variable int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(1)

	// Outer task
	go func() {
		defer wg.Done()

		var innerWg sync.WaitGroup
		innerWg.Add(1)

		// Inner task
		go func() {
			defer innerWg.Done()

			// Thread-private variable for this task
			tp := 1

			// Nested task (task scheduling point)
			var nestedWg sync.WaitGroup
			nestedWg.Add(1)
			go func() {
				defer nestedWg.Done()
				// Empty task - provides scheduling point
			}()
			nestedWg.Wait()

			// Value persists across scheduling point due to constraints
			variable = tp // tp is still 1
		}()

		innerWg.Wait()
	}()

	wg.Wait()
}
