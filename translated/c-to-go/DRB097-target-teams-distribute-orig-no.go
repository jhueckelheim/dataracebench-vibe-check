/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB097-target-teams-distribute-orig-no.c

Description: Use of omp target + teams + distribute + parallel for
GPU offloading with hierarchical parallelism simulation.
*/

package main

import (
	"fmt"
	"sync"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	const len = 2560
	var sum, sum2 float64 = 0.0, 0.0
	a := make([]float64, len)
	b := make([]float64, len)

	// Initialize with some values
	for i := 0; i < len; i++ {
		a[i] = float64(i) / 2.0
		b[i] = float64(i) / 3.0
	}

	// Simulate target + teams + distribute + parallel for
	// In Go, we simulate this with a hierarchical parallel structure
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Simulate teams (10 teams with 256 thread limit)
	numTeams := 10
	teamSize := 256
	chunkSize := 256

	for team := 0; team < numTeams; team++ {
		wg.Add(1)
		go func(teamID int) {
			defer wg.Done()

			// Each team processes different chunks (distribute)
			start := teamID * chunkSize
			if start >= len {
				return
			}

			// Parallel reduction within each team
			var teamSum float64 = 0.0
			var teamWg sync.WaitGroup
			var teamMu sync.Mutex

			threadsPerTeam := min(teamSize, len-start)
			itemsPerThread := min(chunkSize, len-start) / threadsPerTeam
			if itemsPerThread == 0 {
				itemsPerThread = 1
			}

			// Parallel for within team
			for t := 0; t < threadsPerTeam; t++ {
				teamWg.Add(1)
				go func(threadID int) {
					defer teamWg.Done()

					threadStart := start + threadID*itemsPerThread
					threadEnd := min(threadStart+itemsPerThread, min(start+chunkSize, len))

					var localSum float64 = 0.0
					for i := threadStart; i < threadEnd; i++ {
						localSum += a[i] * b[i]
					}

					// Team-level reduction
					teamMu.Lock()
					teamSum += localSum
					teamMu.Unlock()
				}(t)
			}

			teamWg.Wait()

			// Global reduction across teams
			mu.Lock()
			sum += teamSum
			mu.Unlock()
		}(team)
	}

	wg.Wait()

	// CPU reference computation
	var wg2 sync.WaitGroup
	var mu2 sync.Mutex
	numThreads := 4
	itemsPerThread := len / numThreads

	for t := 0; t < numThreads; t++ {
		wg2.Add(1)
		go func(threadID int) {
			defer wg2.Done()

			start := threadID * itemsPerThread
			end := start + itemsPerThread
			if threadID == numThreads-1 {
				end = len
			}

			var localSum float64 = 0.0
			for i := start; i < end; i++ {
				localSum += a[i] * b[i]
			}

			mu2.Lock()
			sum2 += localSum
			mu2.Unlock()
		}(t)
	}

	wg2.Wait()

	fmt.Printf("sum=%f sum2=%f\n", sum, sum2)
}
