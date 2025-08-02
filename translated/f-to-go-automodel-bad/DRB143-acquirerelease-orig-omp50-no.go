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
//No Data Race pair

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var x, y, thrd int32
	var tmp int32
	x = 0

	//$omp parallel num_threads(2) private(thrd) private(tmp)
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(threadID int) {
			defer wg.Done()
			thrd = int32(threadID)
			if thrd == 0 {
				//$omp critical
				var mu sync.Mutex
				mu.Lock()
				x = 10
				mu.Unlock()
				//$omp end critical

				//$omp flush(x)
				atomic.StoreInt32(&x, 10)

				//$omp atomic write
				atomic.StoreInt32(&y, 1)
				//$omp end atomic
			} else {
				tmp = 0
				for tmp == 0 {
					//$omp atomic read acquire ! or seq_cst
					tmp = atomic.LoadInt32(&x)
					//$omp end atomic
				}
				//$omp critical
				fmt.Printf("x = %d\n", x)
				//$omp end critical
			}
		}(i)
	}
	wg.Wait()
	//$omp end parallel
}
