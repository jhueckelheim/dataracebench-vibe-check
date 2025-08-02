//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//For the case of a variable which is referenced within a construct:
//static data member should be shared, unless it is within a threadprivate directive.
//
//Dependence pair: counter@37:5:W vs. counter@37:5:W

package main

import (
	"fmt"
	"sync"
)

// Module DRB087 translated to package-level variables and types
var counter int
var pcounter int

type A struct {
	counter  int
	pcounter int
}

func main() {
	var c A
	c = A{0, 0}

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		counter = counter + 1
		pcounter = pcounter + 1
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d %d\n", counter, pcounter)
}
