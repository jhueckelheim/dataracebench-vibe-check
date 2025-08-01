/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
A single thread will spawn all the tasks. Without if(0) the tasks are deferred and cause data races.

Data Race pairs: var (write vs. write)
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var variable int
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		// Spawn all tasks concurrently (deferred execution)
		var taskWg sync.WaitGroup
		taskWg.Add(10)

		for i := 0; i < 10; i++ {
			go func() {
				defer taskWg.Done()
				variable++ // Data race: concurrent writes to variable
			}()
		}

		taskWg.Wait()
	}()

	wg.Wait()

	// Allow tasks to complete
	time.Sleep(time.Millisecond)

	if variable != 10 {
		fmt.Printf("%d\n", variable)
	}
}
