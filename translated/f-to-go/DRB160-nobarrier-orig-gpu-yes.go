//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Vector addition followed by multiplication involving the same var should have a barrier in between.
//omp distribute directive does not have implicit barrier. This will cause data race.
//Data Race Pair: b[i]@36:23:R vs. b[i]@42:13:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variables (module equivalent)
var a, val int
var b, c, temp [8]int

func main() {
	// Initialize arrays
	for i := 1; i <= 8; i++ {
		b[i-1] = 0
		c[i-1] = 2
		temp[i-1] = 0
	}

	a = 2
	val = 0

	//$omp target map(tofrom:b) map(to:c,temp,a) device(0)
	//$omp teams
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()

	for team := 0; team < numTeams; team++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for i := 1; i <= 100; i++ {
				//$omp distribute
				// NO implicit barrier with distribute
				go func() {
					for j := 1; j <= 8; j++ {
						temp[j-1] = b[j-1] + c[j-1] // Reading b
					}
				}()
				
				//$omp distribute
				// NO barrier between distributes - RACE!
				k := 1 // Fix undefined k
				go func() {
					for j := 8; j >= 1; j -= k-1 { // This will loop infinitely due to k-1=0
						b[j-1] = temp[j-1] * a // RACE: Writing b while first distribute might read
					}
				}()
			}
		}()
	}
	wg.Wait()
	//$omp end teams
	//$omp end target

	for i := 1; i <= 100; i++ {
		val = val + 2
		val = val * 2
	}

	for i := 1; i <= 8; i++ {
		if val != b[i-1] {
			fmt.Printf("%d %d\n", b[i-1], val)
		}
	}
}