/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
The distribute parallel for directive will execute loop using multiple teams.
The loop iterations are distributed across the teams in chunks in round robin fashion.
The omp lock is only guaranteed for a contention group, i.e, within a team.
Data Race Pair, var:W vs. var:W across different teams.
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var variable int = 0
	var wg sync.WaitGroup

	// Simulate multiple teams where locks don't protect across teams
	numTeams := 4
	iterationsPerTeam := 25 // 100/4

	for teamID := 0; teamID < numTeams; teamID++ {
		wg.Add(1)
		go func(team int) {
			defer wg.Done()

			// Each team has its own lock (simulating the issue)
			// In OpenMP teams, locks only work within a team
			var teamLock sync.Mutex

			// Each team processes its chunk with multiple threads
			var teamWg sync.WaitGroup
			threadsPerTeam := 4
			teamWg.Add(threadsPerTeam)

			for threadID := 0; threadID < threadsPerTeam; threadID++ {
				go func(thread int) {
					defer teamWg.Done()

					start := thread * iterationsPerTeam / threadsPerTeam
					end := (thread + 1) * iterationsPerTeam / threadsPerTeam

					for i := start; i < end; i++ {
						// Lock only protects within this team's threads
						// RACE: Different teams access variable concurrently
						teamLock.Lock()
						variable++ // RACE: Not protected from other teams
						teamLock.Unlock()
					}
				}(threadID)
			}

			teamWg.Wait()
		}(teamID)
	}

	wg.Wait()
	fmt.Printf("%d\n", variable)
}
