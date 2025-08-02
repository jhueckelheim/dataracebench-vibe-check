//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//tasks with depend clauses to ensure execution order, no data races.

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var i, j, k int
	i = 0

	//$omp parallel
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp single
		var wgSingle sync.WaitGroup
		wgSingle.Add(3)

		//$omp task depend (out:i)
		go func() {
			defer wgSingle.Done()
			time.Sleep(3 * time.Second)
			i = 1
		}()

		//$omp task depend (in:i)
		go func() {
			defer wgSingle.Done()
			j = i
		}()

		//$omp task depend (in:i)
		go func() {
			defer wgSingle.Done()
			k = i
		}()

		wgSingle.Wait()
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	fmt.Printf("j =%3d  k =%3d\n", j, k)

	if j != 1 && k != 1 {
		fmt.Println("Race Condition")
	}
}
