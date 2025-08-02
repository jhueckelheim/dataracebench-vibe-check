//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// * Cover the implicitly determined rule: In an orphaned task generating construct,
// * formal arguments passed by reference are firstprivate.
// * This requires OpenMP 4.5 to work.
// * Earlier OpenMP does not allow a reference type for a variable within firstprivate().
// * No data race pairs.

package main

import (
	"fmt"
	"sync"
)

// Module DRB100 translated to package-level variables and functions
var a []int

func genTask(i int) {
	//$omp task
	a[i-1] = i + 1
	//$omp end task
}

func main() {
	var i int
	a = make([]int, 100)

	//$omp parallel
	var wgParallel sync.WaitGroup
	wgParallel.Add(1)
	go func() {
		defer wgParallel.Done()
		//$omp single
		var wgSingle sync.WaitGroup
		for i = 1; i <= 100; i++ {
			wgSingle.Add(1)
			go func(i int) {
				defer wgSingle.Done()
				genTask(i)
			}(i)
		}
		wgSingle.Wait()
		//$omp end single
	}()
	wgParallel.Wait()
	//$omp end parallel

	for i = 1; i <= 100; i++ {
		if a[i-1] != i+1 {
			fmt.Printf("warning: a(%d) = %d not expected %d\n", i, a[i-1], i+1)
		}
		//        fmt.Printf("%d %d\n", a[i-1], i+1)
	}
}
