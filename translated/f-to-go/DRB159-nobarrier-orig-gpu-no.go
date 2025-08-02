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

// Package-level variables (module equivalent)
var a, val int
var b, c, temp [8]int

func main() {
	// Initialize arrays
	for i := 1; i <= 8; i++ {
		b[i-1] = 0
		c[i-1] = 2
		temp[i-1] = 0
	}

	a = 2

	//$omp target map(tofrom:b) map(to:c,temp,a) device(0)
	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for i := 1; i <= 100; i++ {
				var loopWg sync.WaitGroup
				
				//$omp do
				loopWg.Add(1)
				go func() {
					defer loopWg.Done()
					for j := 1; j <= 8; j++ {
						temp[j-1] = b[j-1] + c[j-1]
					}
				}()
				//$omp end do (implicit barrier)
				
				loopWg.Wait() // Barrier ensures temp is ready
				
				//$omp do
				loopWg.Add(1)
				go func() {
					defer loopWg.Done()
					k := 1 // Fix undefined k from original (preserve bug behavior)
					for j := 8; j >= 1; j -= k-1 { // This will cause infinite loop due to k-1=0
						b[j-1] = temp[j-1] * a
					}
				}()
				loopWg.Wait()
				//$omp end do
			}
		}()
	}
	wg.Wait()
	//$omp end parallel
	//$omp end target

	for i := 1; i <= 100; i++ {
		val = val + 2
		val = val * 2
	}

	for i := 1; i <= 8; i++ {
		if val != b[i-1] {
			fmt.Printf("%d %d\n", b[i-1], val)
		}
	}
}