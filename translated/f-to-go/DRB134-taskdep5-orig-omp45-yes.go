//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//There is no completion restraint on the second child task. Hence, immediately after the first
//taskwait it is unsafe to access the y variable since the second child task may still be
//executing.
//Data Race at y@34:9:W vs. y@41:23:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func foo() {
	var x, y int
	x = 0
	y = 2

	var task1Wg sync.WaitGroup

	//$omp task depend(inout: x) shared(x)
	task1Wg.Add(1)
	go func() {
		defer task1Wg.Done()
		x = x + 1 // 1st Child Task
	}()
	//$omp end task

	//$omp task depend(in: x) depend(inout: y) shared(x, y)
	// This task depends on x but we don't wait for it before reading y
	go func() {
		task1Wg.Wait() // Wait for x dependency
		y = y - x      // 2nd child task - RACE: we read y before this completes
	}()
	//$omp end task

	//$omp task depend(in: x) if(.FALSE.)
	task1Wg.Wait() // 1st taskwait - waits only for x dependency
	//$omp end task

	fmt.Printf("x= %d\n", x) // Safe - task1 completed
	fmt.Printf("y= %d\n", y) // RACE: task2 might still be modifying y

	//$omp taskwait - but this is after the race
}

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			if threadID == 0 {
				foo()
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel
}
