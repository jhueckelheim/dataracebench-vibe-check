//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//When if() evalutes to true, this program has data races due to true dependence within the loop at 31.
//Data race pair: a[i+1]@32:9:W vs. a[i]@32:18:R

package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	var length, j int
	var u float64
	var a []float64

	length = 100
	a = make([]float64, length)

	// Initialize array
	for i := 1; i <= length; i++ {
		a[i-1] = float64(i)
	}

	rand.Seed(time.Now().UnixNano())
	u = rand.Float64()
	j = int(math.Floor(100 * u))

	// Conditional parallelization based on if clause
	if j%2 == 0 {
		//$omp parallel do if (MOD(j,2)==0) - condition is true, so parallel
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := (length - 1) / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 1; start <= length-1; start += chunkSize {
			end := start + chunkSize - 1
			if end > length-1 {
				end = length - 1
			}
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for i := start; i <= end; i++ {
					a[i] = a[i-1] + 1 // RACE: True dependence when parallel
				}
			}(start, end)
		}
		wg.Wait()
	} else {
		// Sequential execution when condition is false
		for i := 1; i <= length-1; i++ {
			a[i] = a[i-1] + 1 // No race - sequential execution
		}
	}
	//$omp end parallel do

	fmt.Printf("a(50) = %f\n", a[49])

	// deallocate(a)
}
