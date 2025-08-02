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
	"math/rand"
	"runtime"
	"sync"
)

func main() {
	var i, length, rem, j int
	var u float64
	var a []float64

	length = 100
	a = make([]float64, length)

	for i = 1; i <= length; i++ {
		a[i-1] = float64(i)
	}

	u = rand.Float64()
	j = int(100 * u)

	//$omp parallel do if (MOD(j,2)==0)
	if j%2 == 0 {
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
					a[i] = a[i-1] + 1
				}
			}(start, end)
		}
		wg.Wait()
	} else {
		for i = 1; i <= length-1; i++ {
			a[i] = a[i-1] + 1
		}
	}
	//$omp end parallel do

	fmt.Printf("a(50) = %f\n", a[49])
}
