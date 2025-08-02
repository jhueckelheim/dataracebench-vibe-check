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
	"sync"
)

func main() {
	var section_count int
	section_count = 0

	// Force single thread execution (equivalent to omp_set_num_threads(1))

	//$omp parallel
	var wg sync.WaitGroup
	numThreads := 1 // Force single thread

	for threadID := 0; threadID < numThreads; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp sections firstprivate(section_count)
			// Each section gets a private copy of section_count with initial value
			var sectionWg sync.WaitGroup

			// Section 1
			sectionWg.Add(1)
			go func() {
				defer sectionWg.Done()
				section_count_copy := section_count // firstprivate copy
				section_count_copy = section_count_copy + 1
				fmt.Printf("section_count = %8d\n", section_count_copy)
			}()

			// Section 2
			sectionWg.Add(1)
			go func() {
				defer sectionWg.Done()
				section_count_copy := section_count // firstprivate copy
				section_count_copy = section_count_copy + 1
				fmt.Printf("section_count = %8d\n", section_count_copy)
			}()

			sectionWg.Wait()
			//$omp end sections
		}()
	}
	wg.Wait()
	//$omp end parallel
}
