/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file
for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * Data race between non-sibling tasks with declared task dependency fixed by
 * adding a taskwait.
 * Derived from code in https://hal.archives-ouvertes.fr/hal-02177469/document,
 * Listing 1.2
 * No Data Race Pair
 * */

package main

import (
	"fmt"
	"sync"
)

func foo() {
	a := 0
	var wg sync.WaitGroup

	// Simulate parallel single with nested tasks
	wg.Add(2) // Two parent tasks

	// First parent task with nested child task
	go func() {
		defer wg.Done()

		var childWg sync.WaitGroup
		childWg.Add(1)

		// Child task (depend(inout: a))
		go func() {
			defer childWg.Done()
			a++ // Safe: taskwait ensures serialization
		}()

		// Taskwait - wait for child to complete
		childWg.Wait()
	}()

	// Second parent task with nested child task
	go func() {
		defer wg.Done()

		var childWg sync.WaitGroup
		childWg.Add(1)

		// Child task (depend(inout: a))
		go func() {
			defer childWg.Done()
			a++ // Safe: taskwait ensures serialization
		}()

		// Taskwait - wait for child to complete
		childWg.Wait()
	}()

	wg.Wait()
	fmt.Printf("a=%d\n", a)
}

func main() {
	foo()
}
