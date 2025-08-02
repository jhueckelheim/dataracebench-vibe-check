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

func foo() {
	var x, y int
	x = 0
	y = 2

	//$omp task depend(inout: x) shared(x)
	go func() {
		x = x + 1 //1st Child Task
	}()
	//$omp end task

	//$omp task shared(y)
	go func() {
		y = y - 1 //2nd child task
	}()
	//$omp end task

	//$omp task depend(in: x) if(.FALSE.)    //1st taskwait
	//$omp end task

	fmt.Printf("x= %d\n", x)

	//$omp taskwait                          //2nd taskwait

	fmt.Printf("y= %d\n", y)
}

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//$omp single
		foo()
		//$omp end single
	}()
	wg.Wait()
	//$omp end parallel
}
