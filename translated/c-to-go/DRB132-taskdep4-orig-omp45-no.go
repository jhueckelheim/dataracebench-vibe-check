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
		y-- // Modifies y
	}()

	// 1st taskwait - only waits for tasks that depend on x
	firstTaskDone.Wait()

	fmt.Printf("x=%d\n", x)

	// 2nd taskwait - ensures all tasks complete
	wg.Wait()

	fmt.Printf("y=%d\n", y) // Safe: all tasks completed before accessing y
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
