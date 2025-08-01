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
 * Data Race: y (write vs. read)
 */
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

	// 1st Child Task with dependency on x
	go func() {
		defer wg.Done()
		defer firstTaskDone.Done()
		x++ // Task depends on inout: x
	}()

	// 2nd child task (no dependency)
	go func() {
		defer wg.Done()
		y-- // Data race: concurrent write to y
	}()

	// 1st taskwait - only waits for tasks that depend on x
	firstTaskDone.Wait()

	fmt.Printf("x=%d\n", x)
	fmt.Printf("y=%d\n", y) // Data race: read y while 2nd task may still be writing

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
