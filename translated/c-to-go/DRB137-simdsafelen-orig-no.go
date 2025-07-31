/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* The safelen(2) clause guarantees that the vector code is safe for vectors up to 2 (inclusive).
 * In the loop, m can be 2 or more for the correct execution. If the value of m is less than 2,
 * the behavior is undefined. No Data Race in b[i] assignment.
 * */
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	m := 2 // Safe distance for vector operations
	n := 4
	b := [4]int{0, 0, 0, 0}

	var wg sync.WaitGroup

	// Simulate SIMD parallel execution with safe vector length
	// Each iteration can safely run in parallel due to m >= 2 spacing
	for i := m; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// Safe: b[idx] and b[idx-m] don't overlap when m >= 2
			b[idx] = b[idx-m] - 1 // No data race due to sufficient spacing
		}(i)
	}

	wg.Wait()

	fmt.Printf("Expected: -1; Real: %d\n", b[3])
}
