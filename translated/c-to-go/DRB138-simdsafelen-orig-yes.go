/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* The safelen(2) clause guarantees that the vector code is safe for vectors
 * up to 2 (inclusive). In the loop, m can be 2 or more for the correct execution. If the
 * value of m is less than 2, the behavior is undefined.
 * Data Race Pair: b[i] (write) vs. b[i-m] (read)
 * */
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	m := 1 // Unsafe distance for vector operations (less than safelen(2))
	n := 4
	b := [4]int{0, 0, 0, 0}

	var wg sync.WaitGroup

	// Simulate SIMD parallel execution with unsafe vector length
	// Iterations may race when m < 2 (violates safelen constraint)
	for i := m; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// Data race: when m=1, b[2] writes while b[1] may be read by another iteration
			// and b[3] writes while b[2] may be written by another iteration
			b[idx] = b[idx-m] - 1 // Data race: overlapping memory access
		}(i)
	}

	wg.Wait()

	fmt.Printf("Expected: -1; Real: %d\n", b[3])
}
