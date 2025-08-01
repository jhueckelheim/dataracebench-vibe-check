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
 * Data race between non-sibling tasks created from different implicit tasks
 * with declared task dependency
 * Derived from code in https://hal.archives-ouvertes.fr/hal-02177469/document,
 * Listing 1.3
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

	// Simulate parallel with multiple implicit tasks
	numGoroutines := 4
	wg.Add(numGoroutines)

	for t := 0; t < numGoroutines; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Each thread creates a task (depend(inout: a))
			// RACE: Tasks from different implicit tasks can run concurrently
			a++ // RACE: Multiple threads writing to shared variable
		}(t)
	}

	wg.Wait()
	fmt.Printf("a=%d\n", a)
}

func main() {
	foo()
}
