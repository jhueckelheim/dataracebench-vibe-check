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
 * guarantees completion of the second task before y is accessed. Therefore there is no race
 * condition.
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

	// 1st taskwait - waits for task1 (depend(in: x))
	task1Done.Wait()

	// Safe to access x after task1 completion
	fmt.Printf("x=%d\n", x)

	// Task 2 depends on task1 - start after task1 completes
	task2Done.Add(1)
	go func() {
		defer task2Done.Done()
		y = y - x // 2nd child task (now safe to access x)
	}()

	// 2nd taskwait - waits for all remaining tasks
	task2Done.Wait()

	// Safe to access y after task2 completion
	fmt.Printf("y=%d\n", y)
}

func main() {
	// Simulate parallel single
	foo()
}
