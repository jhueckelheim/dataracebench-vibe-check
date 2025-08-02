//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Addition of mutexinoutset dependence type on c, will ensure that line d@36:9 assignment will depend
//on task at Line 29 and line 32. They might execute in any order but not at the same time.
//There is no data race.

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
		var mu sync.Mutex

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

		//$omp task depend(in: a) depend(mutexinoutset: c)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			mu.Lock()
			c = c + a // Task T4
			mu.Unlock()
		}()
		//$omp end task

		//$omp task depend(in: b) depend(mutexinoutset: c)
		wgSingle.Add(1)
		go func() {
			defer wgSingle.Done()
			mu.Lock()
			c = c + b // Task T5
			mu.Unlock()
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
