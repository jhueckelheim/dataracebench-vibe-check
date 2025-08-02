//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//For a variable declared in a scope inside an OpenMP construct:
//* private if the variable has an automatic storage duration
//* shared if the variable has a static storage duration.
//
//Dependence pairs:
//   tmp@38:13:W vs. tmp@38:13:W
//   tmp@38:13:W vs. tmp@39:20:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable simulates static storage (save attribute)
var tmp int // Static storage - shared across threads - RACE CONDITION

func main() {
	var length int
	var a, b []int

	length = 100
	a = make([]int, length)
	b = make([]int, length)

	// Initialize arrays
	for i := 1; i <= length; i++ {
		a[i-1] = i
		b[i-1] = i
	}

	//$omp parallel
	var wg1 sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()

			//$omp do
			chunkSize := length / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = length
			}

			for i := start; i <= end; i++ {
				tmp = a[i-1] + i // RACE: tmp is static (shared) - multiple threads modify
				a[i-1] = tmp     // RACE: Reading shared tmp
			}
			//$omp end do
		}()
	}
	wg1.Wait()
	//$omp end parallel

	//$omp parallel
	var wg2 sync.WaitGroup

	for threadID := 0; threadID < numCPU; threadID++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()

			//$omp do
			chunkSize := length / numCPU
			start := threadID*chunkSize + 1
			end := start + chunkSize - 1
			if threadID == numCPU-1 {
				end = length
			}

			for i := start; i <= end; i++ {
				tmp2 := b[i-1] + i // No race - tmp2 has automatic storage (private to each goroutine)
				b[i-1] = tmp2      // No race - using private tmp2
			}
			//$omp end do
		}()
	}
	wg2.Wait()
	//$omp end parallel

	fmt.Printf("%3d   %3d\n", a[49], b[49])

	// deallocate(a,b)
}
