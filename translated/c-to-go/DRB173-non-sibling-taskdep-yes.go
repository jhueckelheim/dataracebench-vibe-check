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
 * Data race between non-sibling tasks with declared task dependency
 * Derived from code in https://hal.archives-ouvertes.fr/hal-02177469/document,
 * Listing 1.1
 * Data Race Pair, a:W vs. a:W
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
			// RACE: Non-sibling tasks can execute concurrently
			a++ // RACE: Writing to shared variable
		}()

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
			// RACE: Non-sibling tasks can execute concurrently
			a++ // RACE: Writing to shared variable
		}()

		childWg.Wait()
	}()

	wg.Wait()
	fmt.Printf("a=%d\n", a)
}

func main() {
	foo()
}
