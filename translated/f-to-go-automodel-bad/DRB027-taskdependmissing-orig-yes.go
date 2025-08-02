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
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp single
		var wgSingle sync.WaitGroup
		wgSingle.Add(2)

		//$omp task
		go func() {
			defer wgSingle.Done()
			i = 1
		}()
		//$omp end task

		//$omp task
		go func() {
			defer wgSingle.Done()
			i = 2
		}()
		//$omp end task

		wgSingle.Wait()
		//$omp end single
	}()
	//$omp end parallel

	wgParallel.Wait()

	fmt.Printf("i=%3d\n", i)
}
