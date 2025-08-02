//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//The thread encountering the taskwait directive at line 22 only waits for its child task
//(line 14-21) to complete. It does not wait for its descendant tasks (line 16-19). Data Race pairs, sum@36:13:W vs. sum@36:13:W

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var a, psum []int
	var sum, i int

	a = make([]int, 4)
	psum = make([]int, 4)

	//$omp parallel num_threads(2)
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp do schedule(dynamic, 1)
		var wgDo sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := 4 / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 1; start <= 4; start += chunkSize {
			end := start + chunkSize - 1
			if end > 4 {
				end = 4
			}
			wgDo.Add(1)
			go func(start, end int) {
				defer wgDo.Done()
				for i := start; i <= end; i++ {
					a[i-1] = i
				}
			}(start, end)
		}
		wgDo.Wait()
		//$omp end do

		//$omp single
		var wgSingle sync.WaitGroup
		//$omp task
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			//$omp task
			go func() {
				psum[1] = a[2] + a[3]
			}()
			//$omp end task
			psum[0] = a[0] + a[1]
		}()
		//$omp end task
		//$omp taskwait
		wgSingle.Wait()
		sum = psum[1] + psum[0]
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	fmt.Printf("sum = %d\n", sum)
}
