//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Input dependence race: example from OMPRacer: A Scalable and Precise Static Race
// Detector for OpenMP Programs
// Data Race Pair, a(1)@63:26:W vs. a(i)@62:9:W

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func loadFromInput(a []int, N int) {
	for i := 1; i <= N; i++ {
		a[i-1] = i
	}
}

func main() {
	var N int
	var a []int

	N = 100

	argCount := len(os.Args) - 1
	if argCount == 0 {
		fmt.Printf("No command line arguments provided.\n")
	}

	if argCount >= 1 {
		var rdErr error
		N, rdErr = strconv.Atoi(os.Args[1])
		if rdErr != nil {
			fmt.Printf("Error, invalid integer value.\n")
		}
	}

	a = make([]int, N)

	loadFromInput(a, N)

	//$omp parallel do shared(a)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := N / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= N; start += chunkSize {
		end := start + chunkSize - 1
		if end > N {
			end = N
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				a[i-1] = i                      // Writing a[i]
				if N > 10000 {
					a[0] = 1 // RACE: Multiple threads may write to a[1] (a[0] in Go)
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}