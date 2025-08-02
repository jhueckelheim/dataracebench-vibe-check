//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A variable is declared inside a function called within a parallel region.
//The variable should be shared if it uses static storage.
//
//Data race pair: i@19:7:W vs. i@19:7:W

package main

import (
	"fmt"
	"sync"
)

// Module global_foo translated to package-level functions
var i int // Static storage equivalent

func foo() {
	i = i + 1
	fmt.Printf("%d\n", i)
}

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		foo()
	}()
	wg.Wait()
	//$omp end parallel
}
