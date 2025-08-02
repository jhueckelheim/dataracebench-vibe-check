//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Classic PI calculation using reduction. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var x, intervalWidth, pi float64
	var numSteps int64

	pi = 0.0
	numSteps = 2000000000
	intervalWidth = 1.0 / float64(numSteps)

	//$omp parallel do reduction(+:pi) private(x)
	var wg sync.WaitGroup
	var piMutex sync.Mutex
	numCPU := runtime.NumCPU()
	chunkSize := int(numSteps) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= int(numSteps); start += chunkSize {
		end := start + chunkSize - 1
		if end > int(numSteps) {
			end = int(numSteps)
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localPi := 0.0
			for i := start; i <= end; i++ {
				x = (float64(i) + 0.5) * intervalWidth
				localPi = localPi + 1.0/(x*x+1.0)
			}
			piMutex.Lock()
			pi += localPi
			piMutex.Unlock()
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	pi = pi * 4.0 * intervalWidth
	fmt.Printf("PI =%24.20f\n", pi)
}
