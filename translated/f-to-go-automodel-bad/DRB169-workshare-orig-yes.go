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
	chunkSize := 1 / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 0; start < 1; start += chunkSize {
		end := start + chunkSize - 1
		if end >= 1 {
			end = 0
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			//$omp workshare
			AA = BB
			AA = AA + CC
			//$omp end workshare nowait

			//$omp workshare
			res = AA * 2
			//$omp end workshare
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel

	if res != 6 {
		fmt.Printf("%d\n", res)
	}
}
