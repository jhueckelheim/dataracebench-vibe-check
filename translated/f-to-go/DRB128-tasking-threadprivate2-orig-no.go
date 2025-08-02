//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The scheduling constraints prohibit a thread in the team from executing
//a new task that modifies tp while another such task region tied to
//the same thread is suspended. Therefore, the value written will
//persist across the task scheduling point.
//No Data Race

package main

import (
	"fmt"
	"sync"
)

// Package-level variables (module equivalent)
var tp, variable int // tp is threadprivate in original

func foo() {
	var wg sync.WaitGroup

	//$omp task
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Scheduling constraints ensure proper ordering
		localTp := tp

		//$omp task
		localTp = 1
		//$omp task
		//$omp end task (empty task)
		variable = localTp // No race - scheduling constraints maintained
		//$omp end task
	}()
	//$omp end task

	wg.Wait()
}

func main() {
	foo()
	fmt.Printf("%d\n", variable)
}
