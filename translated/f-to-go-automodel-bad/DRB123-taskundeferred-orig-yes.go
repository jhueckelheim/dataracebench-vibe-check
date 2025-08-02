//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A single thread will spawn all the tasks. Add if(0) to avoid the data race, undeferring the tasks.
//Data Race Pairs, var@21:9:W vs. var@21:9:W

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
			//$omp task shared(var)
			wg.Add(1)
			go func() {
				defer wg.Done()
				var1 = var1 + 1
			}()
			//$omp end task
		}
	}()
	wg.Wait()
	//$omp end parallel sections

	fmt.Printf("var =%8d\n", var1)
}
