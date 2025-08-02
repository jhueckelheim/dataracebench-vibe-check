//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Number of threads is empirical: We need enough threads so that
//the reduction is really performed hierarchically in the barrier!
//There is no data race.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var var1, i, sum1, sum2 int

	var1 = 0
	sum1 = 0
	sum2 = 0

	//$omp parallel reduction(+: var)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()
	chunkSize := 5 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	//$omp do schedule(static) reduction(+: sum1)
	for start := 1; start <= 5; start += chunkSize {
		end := start + chunkSize - 1
		if end > 5 {
			end = 5
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localSum1 := 0
			for i := start; i <= end; i++ {
				localSum1 = localSum1 + i
			}
			mu.Lock()
			sum1 = sum1 + localSum1
			mu.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end do

	//$omp do schedule(static) reduction(+: sum2)
	for start := 1; start <= 5; start += chunkSize {
		end := start + chunkSize - 1
		if end > 5 {
			end = 5
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localSum2 := 0
			for i := start; i <= end; i++ {
				localSum2 = localSum2 + i
			}
			mu.Lock()
			sum2 = sum2 + localSum2
			mu.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end do

	var1 = sum1 + sum2
	//$omp end parallel

	fmt.Printf("var =%8d\n", var1)
}
