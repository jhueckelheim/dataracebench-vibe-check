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

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func main() {
	var length int64
	var sum, sum2 float64
	var a, b []float64

	length = 2560
	sum = 0.0
	sum2 = 0.0

	a = make([]float64, length)
	b = make([]float64, length)

	// Initialize arrays
	for i := int64(1); i <= length; i++ {
		a[i-1] = float64(i) / 2.0
		b[i-1] = float64(i) / 3.0
	}

	//$omp target map(to: a(0:len), b(0:len)) map(tofrom: sum)
	//$omp teams num_teams(10) thread_limit(256) reduction (+:sum)
	//$omp distribute
	// Simulate target teams distribute with goroutines
	var wg1 sync.WaitGroup
	var mu1 sync.Mutex
	numTeams := 10

	for team := 0; team < numTeams; team++ {
		wg1.Add(1)
		go func(team int) {
			defer wg1.Done()
			teamSum := 0.0

			// Distribute iterations across teams
			for i2 := int64(team*256 + 1); i2 <= length; i2 += int64(numTeams * 256) {
				//$omp parallel do reduction (+:sum)
				var wg2 sync.WaitGroup
				var mu2 sync.Mutex
				localSum := 0.0
				threadLimit := 256
				numCPU := runtime.NumCPU()
				if numCPU > threadLimit {
					numCPU = threadLimit
				}

				endLoop := min(i2+255, length)
				chunkSize := int(endLoop-i2+1) / numCPU
				if chunkSize < 1 {
					chunkSize = 1
				}

				for start := i2 + 1; start <= endLoop; start += int64(chunkSize) {
					end := start + int64(chunkSize) - 1
					if end > endLoop {
						end = endLoop
					}
					wg2.Add(1)
					go func(start, end int64) {
						defer wg2.Done()
						threadSum := 0.0
						for i := start; i <= end; i++ {
							threadSum += a[i-1] * b[i-1] // No race - proper partitioning
						}
						mu2.Lock()
						localSum += threadSum
						mu2.Unlock()
					}(start, end)
				}
				wg2.Wait()
				teamSum += localSum
			}

			mu1.Lock()
			sum += teamSum
			mu1.Unlock()
		}(team)
	}
	wg1.Wait()
	//$omp end distribute
	//$omp end teams
	//$omp end target

	//$omp parallel do reduction (+:sum2)
	var wg3 sync.WaitGroup
	var mu3 sync.Mutex
	numCPU := runtime.NumCPU()
	chunkSize := int(length) / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := int64(1); start <= length; start += int64(chunkSize) {
		end := start + int64(chunkSize) - 1
		if end > length {
			end = length
		}
		wg3.Add(1)
		go func(start, end int64) {
			defer wg3.Done()
			localSum := 0.0
			for i := start; i <= end; i++ {
				localSum += a[i-1] * b[i-1]
			}
			mu3.Lock()
			sum2 += localSum
			mu3.Unlock()
		}(start, end)
	}
	wg3.Wait()
	//$omp end parallel do

	fmt.Printf("sum = %d; sum2 = %d\n", int(sum), int(sum2))

	// deallocate(a,b)
}
