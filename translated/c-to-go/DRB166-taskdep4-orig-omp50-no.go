/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* The second taskwait ensures that the second child task has completed; hence it is safe to access
 * the y variable in the following print statement.
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
		x++ // 1st Child Task
	}()

	// Task 2 with no dependencies - runs independently
	task2Done.Add(1)
	go func() {
		defer task2Done.Done()
		y-- // 2nd child task
	}()

	// 1st taskwait - only waits for task1 (depend(in: x))
	task1Done.Wait()

	// Safe to access x after task1 completion
	fmt.Printf("x=%d\n", x)

	// 2nd taskwait - waits for all remaining tasks
	task2Done.Wait()

	// Safe to access y after task2 completion
	fmt.Printf("y=%d\n", y)
}

func main() {
	// Simulate parallel single
	foo()
}
