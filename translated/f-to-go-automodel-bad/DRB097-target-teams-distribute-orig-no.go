//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//use of omp target + teams + distribute + parallel for. No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var i, i2, length, lLimit, tmp int64
	var sum, sum2 float64
	var a, b []float64

	length = 2560
	sum = 0.0
	sum2 = 0.0

	a = make([]float64, length)
	b = make([]float64, length)

	for i = 1; i <= length; i++ {
		a[i-1] = float64(i) / 2.0
		b[i-1] = float64(i) / 3.0
	}

	//$omp target map(to: a(0:len), b(0:len)) map(tofrom: sum)
	//$omp teams num_teams(10) thread_limit(256) reduction (+:sum)
	var mu sync.Mutex
	//$omp distribute
	for i2 = 1; i2 <= length; i2 += 256 {
		//$omp parallel do reduction (+:sum)
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		end := min(i2+256, length)
		chunkSize := (end - i2) / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := i2 + 1; start <= end; start += chunkSize {
			chunkEnd := start + chunkSize - 1
			if chunkEnd > end {
				chunkEnd = end
			}
			wg.Add(1)
			go func(start, end int64) {
				defer wg.Done()
				localSum := 0.0
				for i := start; i <= end; i++ {
					localSum = localSum + a[i-1]*b[i-1]
				}
				mu.Lock()
				sum = sum + localSum
				mu.Unlock()
			}(start, chunkEnd)
		}
		wg.Wait()
		//$omp end parallel do
	}
	//$omp end distribute
	//$omp end teams
	//$omp end target

	//$omp parallel do reduction (+:sum2)
	var wg2 sync.WaitGroup
	var mu2 sync.Mutex
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
		wg2.Add(1)
		go func(start, end int64) {
			defer wg2.Done()
			localSum := 0.0
			for i := start; i <= end; i++ {
				localSum = localSum + a[i-1]*b[i-1]
			}
			mu2.Lock()
			sum2 = sum2 + localSum
			mu2.Unlock()
		}(start, end)
	}
	wg2.Wait()
	//$omp end parallel do

	fmt.Printf("sum = %d; sum2 = %d\n", int(sum), int(sum2))
}
