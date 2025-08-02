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
	"sync"
)

func main() {
	var i int
	i = 1

	//$OMP PARALLEL SECTIONS
	var wgSections sync.WaitGroup
	wgSections.Add(1)
	//$OMP SECTION
	go func() {
		defer wgSections.Done()
		//$OMP CRITICAL (NAME)
		var mu sync.Mutex
		mu.Lock()
		//$OMP PARALLEL
		var wgParallel sync.WaitGroup
		wgParallel.Add(1)
		go func() {
			defer wgParallel.Done()
			//$OMP SINGLE
			i = i + 1
			//$OMP END SINGLE
		}()
		wgParallel.Wait()
		//$OMP END PARALLEL
		mu.Unlock()
		//$OMP END CRITICAL (NAME)
	}()
	//$OMP END PARALLEL SECTIONS
	wgSections.Wait()

	fmt.Printf("i = %8d\n", i)
}
