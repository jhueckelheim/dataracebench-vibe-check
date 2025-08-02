//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two tasks with a lock synchronization to ensure execution order. No data race pairs.

package main

import (
	"fmt"
	"sync"
)

func main() {
	var i int
	var mu sync.Mutex

	i = 0

	//$omp parallel sections
	var wg sync.WaitGroup
	wg.Add(2)

	//$omp section
	go func() {
		defer wg.Done()
		mu.Lock()
		i = i + 1
		mu.Unlock()
	}()

	//$omp section
	go func() {
		defer wg.Done()
		mu.Lock()
		i = i + 2
		mu.Unlock()
	}()

	wg.Wait()
	//$omp end parallel sections

	fmt.Printf("I =%3d\n", i)
}
