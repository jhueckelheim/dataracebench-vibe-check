//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//No data race. The data environment of the task is created according to the
//data-sharing attribute clauses, here at line 21:27 it is var. Hence, var is
//modified 10 times, resulting to the value 10.

package main

import (
	"fmt"
	"sync"
)

func main() {
	var var1, i int
	var1 = 0

	//$omp parallel sections
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i = 1; i <= 10; i++ {
			//$omp task shared(var) if(.FALSE.)
			var1 = var1 + 1
			//$omp end task
		}
	}()
	wg.Wait()
	//$omp end parallel sections

	fmt.Printf("var =%8d\n", var1)
}
