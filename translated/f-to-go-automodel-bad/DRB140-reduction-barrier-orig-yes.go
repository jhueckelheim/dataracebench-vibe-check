//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The assignment to a@21:9 is  not synchronized with the update of a@24:30 as a result of the
//reduction computation in the for loop.
//Data Race pair: a@21:9:W vs. a@24:30:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var a, i int

	//$omp parallel shared(a) private(i)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()
	for j := 0; j < numCPU; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//$omp master
			if j == 0 {
				a = 0
			}
			//$omp end master

			//$omp do reduction(+:a)
			chunkSize := 10 / numCPU
			if chunkSize < 1 {
				chunkSize = 1
			}
			localA := 0
			for start := 1; start <= 10; start += chunkSize {
				end := start + chunkSize - 1
				if end > 10 {
					end = 10
				}
				for i := start; i <= end; i++ {
					localA = localA + i
				}
			}
			mu.Lock()
			a = a + localA
			mu.Unlock()
			//$omp end do

			//$omp single
			fmt.Printf("Sum is %d\n", a)
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel
}
