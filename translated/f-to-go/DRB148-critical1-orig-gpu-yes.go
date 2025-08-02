//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Though we have used critical directive to ensure that additions across teams are not overlapped.
//Critical only synchronizes within a team. There is a data race pair.
//Data Race pairs, var@24:9:W vs. var@24:15:R

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
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()

	for team := 0; team < numTeams; team++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Each team has its own critical section
			var teamMutex sync.Mutex
			for i := 1; i <= 100; i++ {
				//$omp critical(addlock)
				teamMutex.Lock() // RACE: Critical only within team, not across teams
				variable = variable + 1 // RACE: Multiple teams modify without global sync
				teamMutex.Unlock()
				//$omp end critical(addlock)
			}
		}()
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", variable)
}