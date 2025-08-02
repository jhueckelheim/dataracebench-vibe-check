//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Addition of mutexinoutset dependence type on c, will ensure that line d@36:9 assignment will depend
//on task at Line 29 and line 32. They might execute in any order but not at the same time.
//There is no data race.

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
				var cMutex sync.Mutex // mutexinoutset equivalent for c

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

				var task4Wg, task5Wg sync.WaitGroup

				//$omp task depend(in: a) depend(mutexinoutset: c)
				task4Wg.Add(1)
				go func() {
					defer task4Wg.Done()
					task2Wg.Wait() // Wait for a
					cMutex.Lock()  // mutexinoutset: exclusive access to c
					c = c + a      // Task T4
					cMutex.Unlock()
				}()

				//$omp task depend(in: b) depend(mutexinoutset: c)
				task5Wg.Add(1)
				go func() {
					defer task5Wg.Done()
					task3Wg.Wait() // Wait for b
					cMutex.Lock()  // mutexinoutset: exclusive access to c
					c = c + b      // Task T5
					cMutex.Unlock()
				}()

				//$omp task depend(in: c)
				go func() {
					task1Wg.Wait() // Wait for initial c
					task4Wg.Wait() // Wait for c modifications
					task5Wg.Wait()
					d = c // Task T6 - no race due to proper dependencies
				}()
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d\n", d)
}
