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
 * second taskwait, there is a race condition. Data Race Pair: y (write vs. read)
 * */
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func foo() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	x := 0
	y := 2

	var wg sync.WaitGroup
	var firstTaskDone sync.WaitGroup

	wg.Add(2)
	firstTaskDone.Add(1)

	// 1st child task
	go func() {
		defer wg.Done()
		defer firstTaskDone.Done()
		x++ // Task modifies x
	}()

	// 2nd child task - depends on x (serialized with first task)
	go func() {
		defer wg.Done()

		// Wait for first task to complete (dependency on x)
		firstTaskDone.Wait()

		y -= x // Data race: concurrent write to y
	}()

	// 1st taskwait - waits for first task (dependency on x)
	firstTaskDone.Wait()

	fmt.Printf("x=%d\n", x)
	fmt.Printf("y=%d\n", y) // Data race: read y while second task may still be writing

	// 2nd taskwait
	wg.Wait()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		foo()
	}()

	wg.Wait()
}
