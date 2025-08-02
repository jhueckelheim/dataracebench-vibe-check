//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//For the case of a variable which is referenced within a construct:
//objects with dynamic storage duration should be shared.
//Putting it within a threadprivate directive may cause seg fault
//since threadprivate copies are not allocated.
//
//Dependence pair: *counter@25:5:W vs. *counter@25:5:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter *int

	counter = new(int)

	*counter = 0

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		*counter = *counter + 1
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d\n", *counter)
}
