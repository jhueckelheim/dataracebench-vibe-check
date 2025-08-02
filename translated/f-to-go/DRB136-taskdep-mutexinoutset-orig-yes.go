//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Due to the missing mutexinoutset dependence type on c, these tasks will execute in any
//order leading to the data race at line 35. Data Race Pair, d@35:9:W vs. d@35:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var a, b, c, d int

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			if threadID == 0 {
				var task1Wg, task2Wg, task3Wg sync.WaitGroup

				//$omp task depend(out: c)
				task1Wg.Add(1)
				go func() {
					defer task1Wg.Done()
					c = 1 // Task T1
				}()

				//$omp task depend(out: a)
				task2Wg.Add(1)
				go func() {
					defer task2Wg.Done()
					a = 2 // Task T2
				}()

				//$omp task depend(out: b)
				task3Wg.Add(1)
				go func() {
					defer task3Wg.Done()
					b = 3 // Task T3
				}()

				//$omp task depend(in: a)
				// MISSING: mutexinoutset dependency on c
				go func() {
					task2Wg.Wait() // Wait for a
					c = c + a      // Task T4 - RACE: modifying c concurrently
				}()

				//$omp task depend(in: b)
				// MISSING: mutexinoutset dependency on c
				go func() {
					task3Wg.Wait() // Wait for b
					c = c + b      // Task T5 - RACE: modifying c concurrently
				}()

				//$omp task depend(in: c)
				go func() {
					// Should wait for all c modifications but dependencies are missing
					d = c // Task T6 - RACE: reading c while T4/T5 might modify it
				}()
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d\n", d)
}
