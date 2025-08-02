//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Depend clause at line 29 and 33 will ensure that there is no data race.

package main

import (
	"fmt"
	"sync"
)

var a, i int
var x [64]int
var y [64]int

func main() {
	for i = 1; i <= 64; i++ {
		x[i-1] = 0
		y[i-1] = 3
	}

	a = 5

	//$omp target map(to:y,a) map(tofrom:x) device(0)
	var wg sync.WaitGroup
	for i = 1; i <= 64; i++ {
		//$omp task depend(inout:x(i))
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			x[idx] = a * x[idx]
		}(i - 1)

		//$omp task depend(inout:x(i))
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			x[idx] = x[idx] + y[idx]
		}(i - 1)
	}
	wg.Wait()
	//$omp end target

	for i = 1; i <= 64; i++ {
		if x[i-1] != 3 {
			fmt.Printf("%d\n", x[i-1])
		}
	}

	//$omp taskwait
}
