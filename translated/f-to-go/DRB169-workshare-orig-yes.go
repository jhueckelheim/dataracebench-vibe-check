//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The workshare construct is only available in Fortran. The workshare spreads work across the threads 
//executing the parallel. There is an implicit barrier. The nowait nullifies this barrier and hence
//there is a race at line:29 due to nowait at line:26. Data Race Pairs, AA@25:9:W vs. AA@29:15:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var AA, BB, CC, res int

	BB = 1
	CC = 2

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			//$omp workshare
			// Workshare distributes work across threads
			if threadID == 0 { // Simulate work distribution
				AA = BB
				AA = AA + CC
			}
			//$omp end workshare nowait
			// NOWAIT: No barrier here - causes race!

			//$omp workshare
			if threadID == 0 { // Second workshare
				res = AA * 2 // RACE: Reading AA while first workshare might still be writing
			}
			//$omp end workshare
		}()
	}
	wg.Wait()
	//$omp end parallel

	if res != 6 {
		fmt.Printf("%d\n", res)
	}
}