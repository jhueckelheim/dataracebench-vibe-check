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
	var variable int

	variable = 0

	//$omp parallel reduction(+: var)
	var wg sync.WaitGroup
	var mu sync.Mutex
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Each thread has private copies for reductions
			localSum1 := 0
			localSum2 := 0
			localVar := 0

			//$omp do schedule(static) reduction(+: sum1)
			for i := 1; i <= 5; i++ {
				localSum1 = localSum1 + i // Private sum1
			}
			//$omp end do

			//$omp do schedule(static) reduction(+: sum2)
			for i := 1; i <= 5; i++ {
				localSum2 = localSum2 + i // Private sum2
			}
			//$omp end do

			localVar = localSum1 + localSum2 // Private var

			// Reduction operations
			mu.Lock()
			variable += localVar // No race - proper reduction
			mu.Unlock()
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("var = %8d\n", variable)
}
