//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The distribute parallel do directive at line 22 will execute loop using multiple teams.
//The loop iterations are distributed across the teams in chunks in round robin fashion.
//The omp lock is only guaranteed for a contention group, i.e, within a team.
//Data Race Pair, var@25:9:W vs. var@25:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var variable int
	var lck sync.Mutex

	//$omp target map(tofrom:var) device(0)
	//$omp teams distribute parallel do
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()

	for team := 0; team < numTeams; team++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 1; i <= 10; i++ {
				// Each team has its own lock scope
				lck.Lock() // RACE: Lock only works within team, not across teams
				variable = variable + 1 // RACE: Multiple teams modify without global sync
				lck.Unlock()
			}
		}()
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	fmt.Printf("%d\n", variable)
}