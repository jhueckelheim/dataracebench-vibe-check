//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This example is based on fpriv_sections.1.c OpenMP Examples 5.0.0
//The section construct modifies the value of section_count which breaks the independence of the
//section constructs. If the same thread executes both the section one will print 1 and the other
//will print 2. For a same thread execution, there is no data race.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var sectionCount int

	sectionCount = 0

	runtime.GOMAXPROCS(1)

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//$omp sections firstprivate( section_count )
		var wgSections sync.WaitGroup
		wgSections.Add(2)
		//$omp section
		go func() {
			defer wgSections.Done()
			localSectionCount := sectionCount // firstprivate equivalent
			localSectionCount = localSectionCount + 1
			fmt.Printf("section_count =%8d\n", localSectionCount)
		}()

		//$omp section
		go func() {
			defer wgSections.Done()
			localSectionCount := sectionCount // firstprivate equivalent
			localSectionCount = localSectionCount + 1
			fmt.Printf("section_count =%8d\n", localSectionCount)
		}()
		//$omp end sections
		wgSections.Wait()
	}()
	wg.Wait()
	//$omp end parallel
}
