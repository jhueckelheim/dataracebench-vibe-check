//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Vector addition followed by multiplication involving the same var should have a barrier in between.
//Here we have an implicit barrier after parallel for regions. No data race pair.

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

	//$omp target map(tofrom:b) map(to:c,temp,a) device(0)
	//$omp parallel
	var wg sync.WaitGroup
	for i = 1; i <= 100; i++ {
		//$omp do
		var wgDo1 sync.WaitGroup
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
			wgDo1.Add(1)
			go func(start, end int) {
				defer wgDo1.Done()
				for j := start; j <= end; j++ {
					temp[j-1] = b[j-1] + c[j-1]
				}
			}(start, end)
		}
		wgDo1.Wait()
		//$omp end do

		//$omp do
		var wgDo2 sync.WaitGroup
		chunkSize = 8 / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 8; start >= 1; start -= chunkSize {
			end := start - chunkSize + 1
			if end < 1 {
				end = 1
			}
			wgDo2.Add(1)
			go func(start, end int) {
				defer wgDo2.Done()
				for j := start; j >= end; j-- {
					b[j-1] = temp[j-1] * a
				}
			}(start, end)
		}
		wgDo2.Wait()
		//$omp end do
	}
	//$omp end parallel
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
