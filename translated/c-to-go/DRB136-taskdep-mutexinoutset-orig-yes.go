/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* Due to the missing mutexinoutset dependence type on c, these tasks will execute in any
 * order leading to the data race. Data Race Pair: c (write vs. write)
 * */
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var a, b, c, d int
	var wg sync.WaitGroup

	var taskA, taskB, taskC sync.WaitGroup
	taskA.Add(1)
	taskB.Add(1)
	taskC.Add(1)

	wg.Add(5)

	// Task: c = 1
	go func() {
		defer wg.Done()

		c = 1 // No protection on c
		taskC.Done()
	}()

	// Task: a = 2
	go func() {
		defer wg.Done()

		a = 2
		taskA.Done()
	}()

	// Task: b = 3
	go func() {
		defer wg.Done()

		b = 3
		taskB.Done()
	}()

	// Task: c += a (depends on a, but missing mutexinoutset on c)
	go func() {
		defer wg.Done()

		taskA.Wait() // Wait for a to be ready

		// Missing mutual exclusion on c - causes data race
		c += a // Data race: concurrent write to c
	}()

	// Task: c += b (depends on b, but missing mutexinoutset on c)
	go func() {
		defer wg.Done()

		taskB.Wait() // Wait for b to be ready

		// Missing mutual exclusion on c - causes data race
		c += b // Data race: concurrent write to c
	}()

	// Wait for all tasks to complete
	wg.Wait()

	// Final task: d = c (depends on c)
	d = c // May read inconsistent value due to races above

	fmt.Printf("%d\n", d)
}
