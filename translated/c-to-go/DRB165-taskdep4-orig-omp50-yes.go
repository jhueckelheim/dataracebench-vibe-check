/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * There is no completion restraint on the second child task. Hence, immediately after the first
 * taskwait it is unsafe to access the y variable since the second child task may still be
 * executing.
 * Data Race at y:W vs. y:R
 */

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
		x++ // 1st Child Task
	}()

	// Task 2 with no dependencies - runs independently
	task2Done.Add(1)
	go func() {
		defer task2Done.Done()
		y-- // 2nd child task - RACE: No dependency tracking
	}()

	// 1st taskwait - only waits for task1 (depend(in: x))
	task1Done.Wait()

	// RACE: Accessing y while task2 might still be running
	fmt.Printf("x=%d\n", x)
	fmt.Printf("y=%d\n", y) // RACE: Reading y while task2 may be writing

	// 2nd taskwait - waits for all remaining tasks
	task2Done.Wait()
}

func main() {
	// Simulate parallel single
	foo()
}
