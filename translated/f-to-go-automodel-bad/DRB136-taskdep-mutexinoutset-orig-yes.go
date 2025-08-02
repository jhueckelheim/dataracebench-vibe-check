//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Due to the missing mutexinoutset dependence type on c, these tasks will execute in any
//order leading to the data race at line 35. Data Race Pair, d@35:9:W vs. d@35:9:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var a, b, c, d int

	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//$omp single
		var wgSingle sync.WaitGroup

		//$omp task depend(out: c)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			c = 1 // Task T1
		}()
		//$omp end task

		//$omp task depend(out: a)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			a = 2 // Task T2
		}()
		//$omp end task

		//$omp task depend(out: b)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			b = 3 // Task T3
		}()
		//$omp end task

		//$omp task depend(in: a)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			c = c + a // Task T4
		}()
		//$omp end task

		//$omp task depend(in: b)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			c = c + b // Task T5
		}()
		//$omp end task

		//$omp task depend(in: c)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			d = c // Task T6
		}()
		//$omp end task

		wgSingle.Wait()
		//$omp end single
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("%d\n", d)
}
