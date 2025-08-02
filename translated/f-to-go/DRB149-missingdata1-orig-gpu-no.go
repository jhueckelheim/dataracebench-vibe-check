//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Classic i-k-j matrix multiplication. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var length int
	var a, b, c []int

	length = 100

	a = make([]int, length)
	b = make([]int, length+length*length)
	c = make([]int, length)

	// Initialize arrays
	for i := 1; i <= length; i++ {
		for j := 1; j <= length; j++ {
			b[j-1+(i-1)*length] = 1
		}
		a[i-1] = 1
		c[i-1] = 0
	}

	//$omp target map(to:a,b) map(tofrom:c) device(0)
	//$omp teams distribute parallel do
	var wg sync.WaitGroup
	numTeams := runtime.NumCPU()
	chunkSize := length / numTeams
	if chunkSize < 1 {
		chunkSize = 1
	}

	for team := 0; team < numTeams; team++ {
		start := team*chunkSize + 1
		end := start + chunkSize - 1
		if team == numTeams-1 {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				for j := 1; j <= length; j++ {
					c[i-1] = c[i-1] + a[j-1]*b[j-1+(i-1)*length] // No race - proper partitioning
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end teams distribute parallel do
	//$omp end target

	for i := 1; i <= length; i++ {
		if c[i-1] != length {
			fmt.Printf("%d\n", c[i-1])
		}
	}

	// deallocate(a,b,c)
}