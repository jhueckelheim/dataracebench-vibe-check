//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This example is referred from OpenMP Application Programming Interface 5.0, example tasking.7.c
//A task switch may occur at a task scheduling point. A single thread may execute both of the
//task regions that modify tp. The parts of these task regions in which tp is modified may be
//executed in any order so the resulting value of var can be either 1 or 2.
//There is a Race pair var@24:13 and var@24:13 but no data race.

package main

// No imports needed

// Package-level variables (module equivalent)
var tp, variable int // tp is threadprivate in original

func foo() {
	//$omp task
	go func() {
		// Each task gets its own tp (threadprivate equivalent)
		localTp := tp

		//$omp task
		go func() {
			localTp = 1
			//$omp task
			//$omp end task (empty task)
			variable = localTp // value can be 1 or 2 due to task scheduling
		}()

		localTp = 2
	}()
	//$omp end task
}

func main() {
	foo()
}
