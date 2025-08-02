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
	"sync"
)

func main() {
	//$omp parallel
	//$omp single
	foo()
	//$omp end single
	//$omp end parallel
}

func foo() {
	var x, y int
	x = 0
	y = 2

	var wg sync.WaitGroup
	var wgX sync.WaitGroup

	//$omp task depend(inout: x) shared(x)
	wgX.Add(1)
	go func() {
		defer wgX.Done()
		x = x + 1 //1st Child Task
	}()
	//$omp end task

	//$omp task shared(y)
	wg.Add(1)
	go func() {
		defer wg.Done()
		y = y - 1 //2nd child task
	}()
	//$omp end task

	//$omp taskwait depend(in: x) //1st taskwait
	wgX.Wait()

	fmt.Printf("x= %d\n", x)

	//$omp taskwait //2nd taskwait
	wg.Wait()

	fmt.Printf("y= %d\n", y)
}
