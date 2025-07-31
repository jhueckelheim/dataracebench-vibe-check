/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * This example is based on fpriv_sections.1.c OpenMP Examples 5.0.0
 * The section construct modifies the value of section_count which breaks the independence of the
 * section constructs. If the same thread executes both sections, one will print 1 and the other
 * will print 2. For same thread execution, there is no data race.
 */
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(1) // Force single thread execution like omp_set_num_threads(1)

	sectionCount := 0
	var wg sync.WaitGroup

	wg.Add(2) // Two sections

	// Section 1 - firstprivate copy of sectionCount
	go func() {
		defer wg.Done()

		// Each section gets its own private copy (firstprivate behavior)
		localSectionCount := sectionCount // Private copy initialized with original value
		localSectionCount++
		fmt.Printf("%d\n", localSectionCount)
	}()

	// Section 2 - firstprivate copy of sectionCount
	go func() {
		defer wg.Done()

		// Each section gets its own private copy (firstprivate behavior)
		localSectionCount := sectionCount // Private copy initialized with original value
		localSectionCount++
		fmt.Printf("%d\n", localSectionCount)
	}()

	wg.Wait()
}
