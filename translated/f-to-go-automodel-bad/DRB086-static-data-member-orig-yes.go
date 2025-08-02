//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//For the case of a variable which is not referenced within a construct:
//static data member should be shared, unless it is within a threadprivate directive.
//
//Dependence pair: counter@27:9:W vs. counter@27:9:W

package main

import (
	"fmt"
	"sync"
)

// Module DRB086 translated to package-level variables and types
var counter int
var pcounter int

type A struct {
	counter  int
	pcounter int
}

func foo() {
	counter = counter + 1
	pcounter = pcounter + 1
}

func main() {
	var c A
	c = A{0, 0}

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		foo()
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d %d\n", counter, pcounter)
}
