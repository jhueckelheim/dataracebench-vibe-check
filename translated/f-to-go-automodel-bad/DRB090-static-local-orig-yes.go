//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//For a variable declared in a scope inside an OpenMP construct:
//* private if the variable has an automatic storage duration
//* shared if the variable has a static storage duration.
//
//Dependence pairs:
//   tmp@38:13:W vs. tmp@38:13:W
//   tmp@38:13:W vs. tmp@39:20:R

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, length int
	var a, b []int
	var tmp int // Static storage equivalent
	var tmp2 int

	length = 100
	a = make([]int, length)
	b = make([]int, length)

	for i = 1; i <= length; i++ {
		a[i-1] = i
		b[i-1] = i
	}

	//$omp parallel
	var wg1 sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := length / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	//$omp do
	for start := 1; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg1.Add(1)
		go func(start, end int) {
			defer wg1.Done()
			for i := start; i <= end; i++ {
				tmp = a[i-1] + i
				a[i-1] = tmp
			}
		}(start, end)
	}
	wg1.Wait()
	//$omp end do
	//$omp end parallel

	//$omp parallel
	var wg2 sync.WaitGroup

	//$omp do
	for start := 1; start <= length; start += chunkSize {
		end := start + chunkSize - 1
		if end > length {
			end = length
		}
		wg2.Add(1)
		go func(start, end int) {
			defer wg2.Done()
			for i := start; i <= end; i++ {
				tmp2 = b[i-1] + i
				b[i-1] = tmp2
			}
		}(start, end)
	}
	wg2.Wait()
	//$omp end do
	//$omp end parallel

	fmt.Printf("%3d   %3d\n", a[49], b[49])
}
