/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* The first two tasks are serialized, because a dependence on the first child is produced
 * by x with the in dependence type in the depend clause of the second task. Generating task
 * at the first taskwait only waits for the first child task to complete. The second taskwait
 * guarantees completion of the second task before y is accessed. If we access y before the
 * second taskwait, there is a race condition. Data Race Pair, y:W vs. y:R
 * */

package main

import (
	"fmt"
	"sync"
)

func foo() {
	x := 0
	y := 2

	var task1Done sync.WaitGroup
	var task2Done sync.WaitGroup

	// Task 1 with dependency on x
	task1Done.Add(1)
	go func() {
		defer task1Done.Done()
		x++ // 1st child task
	}()

	// Task 2 depends on task1 (depend(in: x)) and modifies y
	task2Done.Add(1)
	go func() {
		defer task2Done.Done()
		// Wait for task1 to complete before accessing x
		task1Done.Wait()
		y -= x // 2nd child task (serialized after task1)
	}()

	// 1st taskwait - waits for task1 (depend(in: x))
	task1Done.Wait()

	// Safe to access x after task1 completion
	fmt.Printf("x=%d\n", x)

	// RACE: Accessing y before task2 completes
	fmt.Printf("y=%d\n", y) // RACE: Reading y while task2 may be writing

	// 2nd taskwait - waits for all remaining tasks
	task2Done.Wait()
}

func main() {
	// Simulate parallel single
	foo()
}
