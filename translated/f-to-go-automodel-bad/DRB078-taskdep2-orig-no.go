//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two tasks with depend clause to ensure execution order, no data races.
//i is shared for two tasks based on implicit data-sharing attribute rules.

package main

import (
	"fmt"
	"sync"
	"time"
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
			time.Sleep(3 * time.Second)
			i = 3
		}()

		//$omp task depend (out:i)
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
		fmt.Printf("%3d\n", i)
	}
}
