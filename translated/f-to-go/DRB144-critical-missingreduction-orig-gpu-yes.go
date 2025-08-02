//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The increment at line number 22 is critical for the variable
//var@22:13. Therefore, there is a possible Data Race pair var@22:13:W vs. var@22:19:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var variable int
	variable = 0

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do
	// Simulate GPU teams with multiple goroutines
	var wg sync.WaitGroup
	var criticalMutex sync.Mutex
	numTeams := runtime.NumCPU()
	chunkSize := 200 / numTeams

	for team := 0; team < numTeams; team++ {
		start := team*chunkSize + 1
		end := start + chunkSize - 1
		if team == numTeams-1 {
			end = 200
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				//$omp critical
				criticalMutex.Lock()
				variable = variable + 1 // RACE: Critical only within team, not across teams
				criticalMutex.Unlock()
				//$omp end critical
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", variable)
}