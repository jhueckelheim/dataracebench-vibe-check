//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The second taskwait ensures that the second child task has completed; hence it is safe to access
//the y variable in the following print statement.

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

	var task1Wg, allTasksWg sync.WaitGroup
	
	//$omp task depend(inout: x) shared(x)
	task1Wg.Add(1)
	allTasksWg.Add(1)
	go func() {
		defer task1Wg.Done()
		defer allTasksWg.Done()
		x = x + 1 // 1st Child Task
	}()
	//$omp end task

	//$omp task shared(y)
	allTasksWg.Add(1)
	go func() {
		defer allTasksWg.Done()
		y = y - 1 // 2nd child task
	}()
	//$omp end task

	//$omp taskwait depend(in: x)
	task1Wg.Wait() // 1st taskwait - waits for x dependency

	fmt.Printf("x= %d\n", x) // Safe - x dependency satisfied

	//$omp taskwait
	allTasksWg.Wait() // 2nd taskwait - waits for ALL tasks

	fmt.Printf("y= %d\n", y) // No race - all tasks completed
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