/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Data race in vectorizable code
Loop depencency with 64 element offset. Data race present.
Data Race Pairs, a[i + 64]@34:5:W vs. a[i]@34:17:R
*/

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	length := 20000
	if len(os.Args) > 1 {
		if val, err := strconv.Atoi(os.Args[1]); err == nil {
			length = val
		}
	}

	a := make([]float64, length)
	b := make([]float64, length)

	for i := 0; i < length; i++ {
		a[i] = float64(i)
		b[i] = float64(i + 1)
	}

	// Simulate OpenMP parallel for simd with goroutines
	numWorkers := runtime.NumCPU()
	chunkSize := 64 // dynamic schedule with chunk size 64
	var wg sync.WaitGroup

	maxIter := length - 64
	for start := 0; start < maxIter; start += chunkSize * numWorkers {
		for worker := 0; worker < numWorkers; worker++ {
			wg.Add(1)
			go func(workerStart int) {
				defer wg.Done()
				end := workerStart + chunkSize
				if end > maxIter {
					end = maxIter
				}
				for i := workerStart; i < end; i++ {
					// Race condition: writing a[i+64] while reading a[i]
					// Creates sliding window dependency where later iterations read what earlier iterations wrote
					a[i+64] = a[i] + b[i]
				}
			}(start + worker*chunkSize)
		}
		wg.Wait()
	}

	fmt.Printf("a[0]=%f, a[%d]=%f, a[%d]=%f\n", a[0], length/2, a[length/2], length-1, a[length-1])
}
