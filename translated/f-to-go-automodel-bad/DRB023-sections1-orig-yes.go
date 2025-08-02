//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two tasks without synchronization to protect data write, causing data races.
//Data race pair: i@20:5:W vs. i@22:5:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var i int
	i = 0

	//$omp parallel sections
	var wg sync.WaitGroup
	wg.Add(2)

	//$omp section
	go func() {
		defer wg.Done()
		i = 1
	}()
	//$omp section
	go func() {
		defer wg.Done()
		i = 2
	}()
	//$omp end parallel sections

	wg.Wait()

	fmt.Printf("i=%3d\n", i)
}
