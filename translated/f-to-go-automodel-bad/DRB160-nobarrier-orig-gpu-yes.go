//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Vector addition followed by multiplication involving the same var should have a barrier in between.
//omp distribute directive does not have implicit barrier. This will cause data race.
//Data Race Pair: b[i]@36:23:R vs. b[i]@42:13:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

var a, i, j, k, val int
var b [8]int
var c [8]int
var temp [8]int

func main() {
	for i = 1; i <= 8; i++ {
		b[i-1] = 0
		c[i-1] = 2
		temp[i-1] = 0
	}

	a = 2
	val = 0

	//$omp target map(tofrom:b) map(to:c,temp,a) device(0)
	//$omp teams
	var wg sync.WaitGroup
	for i = 1; i <= 100; i++ {
		//$omp distribute
		var wgDist1 sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := 8 / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 1; start <= 8; start += chunkSize {
			end := start + chunkSize - 1
			if end > 8 {
				end = 8
			}
			wgDist1.Add(1)
			go func(start, end int) {
				defer wgDist1.Done()
				for j := start; j <= end; j++ {
					temp[j-1] = b[j-1] + c[j-1]
				}
			}(start, end)
		}
		wgDist1.Wait()
		//$omp end distribute

		//$omp distribute
		var wgDist2 sync.WaitGroup
		chunkSize = 8 / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 8; start >= 1; start -= chunkSize {
			end := start - chunkSize + 1
			if end < 1 {
				end = 1
			}
			wgDist2.Add(1)
			go func(start, end int) {
				defer wgDist2.Done()
				for j := start; j >= end; j-- {
					b[j-1] = temp[j-1] * a
				}
			}(start, end)
		}
		wgDist2.Wait()
		//$omp end distribute
	}
	//$omp end teams
	//$omp end target

	for i = 1; i <= 100; i++ {
		val = val + 2
		val = val * 2
	}

	for i = 1; i <= 8; i++ {
		if val != b[i-1] {
			fmt.Printf("%d %d\n", b[i-1], val)
		}
	}
}
