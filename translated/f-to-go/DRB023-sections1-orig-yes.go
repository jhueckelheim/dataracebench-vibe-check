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

	//$omp section
	wg.Add(1)
	go func() {
		defer wg.Done()
		i = 1 // Race condition: concurrent write to shared variable
	}()

	//$omp section
	wg.Add(1)
	go func() {
		defer wg.Done()
		i = 2 // Race condition: concurrent write to shared variable
	}()

	wg.Wait()
	//$omp end parallel sections

	fmt.Printf("i=%3d\n", i)
}
