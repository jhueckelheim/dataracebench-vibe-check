//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two-dimensional array computation:
//Only one loop is associated with omp taskloop.
//The inner loop's loop iteration variable will be shared if it is shared in the enclosing context.
//Data race pairs (we allow multiple ones to preserve the pattern):
//  Write_set = {j@36:20 (implicit step +1)}
//  Read_set = {j@36:20, j@37:35}
//  Any pair from Write_set vs. Write_set  and Write_set vs. Read_set is a data race pair.

//need to run with large thread number and large num of iterations.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Package-level variable (module equivalent)
var a [][]int

func main() {
	var length, j int // j is shared among all tasks - RACE CONDITION
	length = 100

	a = make([][]int, length)
	for i := 0; i < length; i++ {
		a[i] = make([]int, length)
	}

	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//$omp single
			if threadID == 0 { // Only one thread executes single region
				//$omp taskloop
				// Each task processes one outer loop iteration
				var taskWg sync.WaitGroup

				for i := 1; i <= length; i++ {
					taskWg.Add(1)
					go func(i int) { // i is private to each task
						defer taskWg.Done()
						for j = 1; j <= length; j++ { // RACE: j is shared among all tasks!
							a[i-1][j-1] = a[i-1][j-1] + 1
						}
					}(i)
				}
				taskWg.Wait()
				//$omp end taskloop
			}
			//$omp end single
		}()
	}
	wg.Wait()
	//$omp end parallel

	fmt.Printf("a(50,50) = %3d\n", a[49][49])
}
