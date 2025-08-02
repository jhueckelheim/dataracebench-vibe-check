//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Concurrent access of var@22:5 in an intra region. Missing Lock leads to intra region data race.
//Data Race pairs, var@22:13:W vs. var@22:13:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var variable int
	variable = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams num_teams(1)
	//$omp distribute parallel do
	// Single team but no lock protection
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// MISSING: Lock protection
			variable = variable + 1 // RACE: Multiple threads modify without protection
		}()
	}
	wg.Wait()
	//$omp end distribute parallel do
	//$omp end teams
	//$omp end target

	fmt.Printf("%d\n", variable)
}