//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//* This is a program based on a test contributed by Yizi Gu@Rice Univ.
//* Use taskgroup to synchronize two tasks. No data race pairs.

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var result int
	result = 0

	//$omp parallel
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp single
		var wgSingle sync.WaitGroup
		//$omp taskgroup
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			//$omp task
			time.Sleep(3 * time.Second)
			result = 1
			//$omp end task
		}()
		wgSingle.Wait()
		//$omp end taskgroup
		//$omp task
		result = 2
		//$omp end task
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	fmt.Printf("result =%8d\n", result)
}
