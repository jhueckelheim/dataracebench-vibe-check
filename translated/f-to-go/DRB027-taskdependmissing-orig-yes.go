//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two tasks without depend clause to protect data writes.
//i is shared for two tasks based on implicit data-sharing attribute rules.
//Data race pair: i@22:5:W vs. i@25:5:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var i int
	i = 0

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		//$omp single
		var taskWg sync.WaitGroup

		//$omp task
		taskWg.Add(1)
		go func() {
			defer taskWg.Done()
			i = 1 // Race condition: concurrent write to shared variable
		}()

		//$omp task
		taskWg.Add(1)
		go func() {
			defer taskWg.Done()
			i = 2 // Race condition: concurrent write to shared variable
		}()

		taskWg.Wait()
		//$omp end single
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("i=%3d\n", i)
}
