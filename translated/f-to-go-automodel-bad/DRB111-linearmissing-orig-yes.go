//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// *  loop missing the linear clause
// *  Data race pair:  j@37:11:R vs. j@38:9:W
// *                   j@37:18:R vs. j@38:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var length, i, j int
	var a, b, c []float64

	length = 100
	i = 0
	j = 0

	a = make([]float64, length)
	b = make([]float64, length)
	c = make([]float64, length)

	for i = 1; i <= length; i++ {
		a[i-1] = float64(i) / 2.0
		b[i-1] = float64(i) / 3.0
		c[i-1] = float64(i) / 7.0
	}

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := length / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				c[j] = c[j] + a[i-1]*b[i-1]
				j = j + 1
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("c(50) = %f\n", c[49])
}
