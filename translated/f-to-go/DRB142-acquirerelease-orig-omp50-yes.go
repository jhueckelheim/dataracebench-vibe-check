//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The below program will fail to order the write to x on thread 0 before the read from x on thread 1.
//The implicit release flush on exit from the critical region will not synchronize with the acquire
//flush that occurs on the atomic read operation performed by thread 1. This is because implicit
//release flushes that occur on a given construct may only synchronize with implicit acquire flushes
//on a compatible construct (and vice-versa) that internally makes use of the same synchronization
//variable.
//
//Implicit flush must be used after critical construct to avoid data race.
//Data Race pair: x@30:13:W vs. x@30:13:W

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var x, y int64 // Use int64 for atomic operations
	x = 0

	//$omp parallel num_threads(2) private(thrd) private(tmp)
	var wg sync.WaitGroup
	var criticalMutex sync.Mutex
	numThreads := 2

	for threadID := 0; threadID < numThreads; threadID++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			if threadID == 0 {
				//$omp critical
				criticalMutex.Lock()
				x = 10 // RACE: Critical section doesn't synchronize with atomic read
				criticalMutex.Unlock()
				//$omp end critical
				// MISSING: explicit flush(x)

				//$omp atomic write
				atomic.StoreInt64(&y, 1)
				//$omp end atomic
			} else {
				tmp := int64(0)
				for tmp == 0 {
					//$omp atomic read acquire
					tmp = atomic.LoadInt64(&x) // RACE: May not see x=10 due to missing sync
					//$omp end atomic
				}
				//$omp critical
				criticalMutex.Lock()
				fmt.Printf("x = %d\n", x)
				criticalMutex.Unlock()
				//$omp end critical
			}
		}(threadID)
	}
	wg.Wait()
	//$omp end parallel
}
