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
	var lock sync.Mutex
	var i int

	i = 0

	//$omp parallel sections
	var wg sync.WaitGroup

	// Section 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		lock.Lock() // omp_set_lock equivalent
		i = i + 1
		lock.Unlock() // omp_unset_lock equivalent
	}()

	// Section 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		lock.Lock() // omp_set_lock equivalent
		i = i + 2
		lock.Unlock() // omp_unset_lock equivalent
	}()

	wg.Wait()
	//$omp end parallel sections

	// omp_destroy_lock equivalent - handled automatically by Go

	fmt.Printf("I = %3d\n", i)
}
