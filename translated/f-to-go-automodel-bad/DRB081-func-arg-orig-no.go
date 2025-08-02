//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A function argument passed by value should be private inside the function.
//Variable i is read only. No data race pairs.

package main

import (
	"fmt"
	"sync"
)

// Module global translated to package-level functions
func f1(i int) int {
	i = i + 1
	return i
}

func main() {
	var i int
	i = 0

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		f1(i)
	}()
	wg.Wait()
	//$omp end parallel

	if i != 0 {
		fmt.Printf("i =%3d\n", i)
	}
}
