//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A single directive is used to protect a write. No data race pairs.

package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	count = 0

	//$omp parallel shared(count)
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(1)
	go func() {
		defer wg.Done()
		//$omp single
		mu.Lock()
		count = count + 1
		mu.Unlock()
		//$omp end single
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("count =%3d\n", count)
}
