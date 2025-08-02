//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//There is no completion restraint on the second child task. Hence, immediately after the first
//taskwait it is unsafe to access the y variable since the second child task may still be
//executing.
//Data Race at y@34:8:W vs. y@40:23:R

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

	//$omp task shared(y)
	go func() {
		y = y - 1 // 2nd child task - NO dependency tracking
	}()
	//$omp end task

	//$omp taskwait depend(in: x)
	task1Wg.Wait() // 1st taskwait - only waits for x dependency

	fmt.Printf("x= %d\n", x) // Safe - x dependency satisfied
	fmt.Printf("y= %d\n", y) // RACE: y task might still be running

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