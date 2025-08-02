//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Referred from worksharing_critical.1.f
//A single thread executes the one and only section in the sections region, and executes the
//critical region. The same thread encounters the nested parallel region, creates a new team
//of threads, and becomes the master of the new team. One of the threads in the new team enters
//the single region and increments i by 1. At the end of this example i is equal to 2.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i int
	i = 1

	//$OMP PARALLEL SECTIONS
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		//$OMP SECTION
		//$OMP CRITICAL (NAME)
		var criticalMutex sync.Mutex
		criticalMutex.Lock()

		//$OMP PARALLEL
		var nestedWg sync.WaitGroup
		var once sync.Once
		numCPU := runtime.NumCPU()

		for threadID := 0; threadID < numCPU; threadID++ {
			nestedWg.Add(1)
			go func() {
				defer nestedWg.Done()

				//$OMP SINGLE
				once.Do(func() {
					i = i + 1 // No race - single execution within critical
				})
				//$OMP END SINGLE
			}()
		}
		nestedWg.Wait()
		//$OMP END PARALLEL

		criticalMutex.Unlock()
		//$OMP END CRITICAL (NAME)
	}()

	wg.Wait()
	//$OMP END PARALLEL SECTIONS

	fmt.Printf("i = %8d\n", i)
}
