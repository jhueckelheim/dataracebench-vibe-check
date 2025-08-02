//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two tasks with depend clause to ensure execution order:
//i is shared for two tasks based on implicit data-sharing attribute rules. No data race pairs.

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

		//$omp task depend (out:i)
		go func() {
			defer wgSingle.Done()
			i = 1
		}()

		//$omp task depend (in:i)
		go func() {
			defer wgSingle.Done()
			i = 2
		}()

		wgSingle.Wait()
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	if i != 2 {
		fmt.Println("i is not equal to 2")
	}
}
