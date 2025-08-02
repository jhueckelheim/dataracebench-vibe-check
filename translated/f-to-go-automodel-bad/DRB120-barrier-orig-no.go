//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The barrier construct specifies an explicit barrier at the point at which the construct appears.
//Barrier construct at line:27 ensures that there is no data race.

package main

import (
	"fmt"
	"sync"
)

func main() {
	var var1 int

	//$omp parallel shared(var)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//$omp single
		var1 = var1 + 1
		//$omp end single
		//$omp barrier
		//$omp single
		var1 = var1 + 1
		//$omp end single
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("var =%3d\n", var1)
}
