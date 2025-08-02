//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Data race on outLen due to ++ operation.
//Adding private (outLen) can avoid race condition. But it is wrong semantically.
//Data races on outLen also cause output[outLen++] to have data races.
//
//Data race pairs (we allow two pairs to preserve the original code pattern):
//1. outLen@34:9:W vs. outLen@34:9:W
//2. output[]@33:9:W vs. output[]@33:9:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var i, inLen, outLen int
	var input [1000]int
	var output [1000]int

	inLen = 1000
	outLen = 1

	for i = 1; i <= inLen; i++ {
		input[i-1] = i
	}

	//$omp parallel do
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := inLen / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= inLen; start += chunkSize {
		end := start + chunkSize - 1
		if end > inLen {
			end = inLen
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				output[outLen-1] = input[i-1]
				outLen = outLen + 1
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do

	fmt.Printf("output(500)=%3d\n", output[499])
}
