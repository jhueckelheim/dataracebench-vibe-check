/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/* Addition of mutexinoutset dependence type on c, will ensure that assignment will depend
 * on previous tasks. They might execute in any order but not at the same time.
 * There is no data race.
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
	var cMutex sync.Mutex // Mutex for mutexinoutset behavior on c

	var taskA, taskB, taskC sync.WaitGroup
	taskA.Add(1)
	taskB.Add(1)
	taskC.Add(1)

	wg.Add(5) // Only 5 tasks that need to complete

	// Task: c = 1
	go func() {
		defer wg.Done()

		cMutex.Lock()
		c = 1
		cMutex.Unlock()
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

	// Task: c += a (depends on a, mutexinoutset on c)
	go func() {
		defer wg.Done()

		taskA.Wait() // Wait for a to be ready
		taskC.Wait() // Wait for c to be initialized

		cMutex.Lock() // Mutual exclusion on c
		c += a
		cMutex.Unlock()
	}()

	// Task: c += b (depends on b, mutexinoutset on c)
	go func() {
		defer wg.Done()

		taskB.Wait() // Wait for b to be ready
		taskC.Wait() // Wait for c to be initialized

		cMutex.Lock() // Mutual exclusion on c - prevents concurrent access
		c += b
		cMutex.Unlock()
	}()

	// Wait for all tasks to complete
	wg.Wait()

	// Final task: d = c (depends on c)
	cMutex.Lock()
	d = c // Safe: all modifications to c completed
	cMutex.Unlock()

	fmt.Printf("%d\n", d)
}
