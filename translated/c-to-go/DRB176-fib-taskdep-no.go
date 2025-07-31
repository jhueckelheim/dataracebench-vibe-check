/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * Fibonacci code without data race
 * No Data Race Pair
 * */

package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func fib(n int) int {
	if n < 2 {
		return n
	}

	var i, j, s int
	var wg sync.WaitGroup

	// Task 1: compute fib(n-1) with dependency out:i
	wg.Add(1)
	go func() {
		defer wg.Done()
		i = fib(n - 1)
	}()

	// Task 2: compute fib(n-2) with dependency out:j
	wg.Add(1)
	go func() {
		defer wg.Done()
		j = fib(n - 2)
	}()

	// Wait for both i and j to be computed (dependencies in:i,j)
	wg.Wait()

	// Task 3: compute sum with dependencies in:i,j out:s
	var sumWg sync.WaitGroup
	sumWg.Add(1)
	go func() {
		defer sumWg.Done()
		s = i + j
	}()

	// Taskwait - wait for sum computation
	sumWg.Wait()

	return s
}

func main() {
	n := 10
	if len(os.Args) > 1 {
		if arg, err := strconv.Atoi(os.Args[1]); err == nil {
			n = arg
		}
	}

	// Simulate parallel sections with single section
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("fib(%d) = %d\n", n, fib(n))
	}()
	wg.Wait()
}
