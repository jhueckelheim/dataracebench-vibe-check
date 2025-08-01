/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
No data race. The tasks are executed immediately (undeferred) due to if(0) condition.
Hence, var is modified 10 times sequentially, resulting in the value 10.
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var variable int
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		// Sequential execution - tasks are undeferred (equivalent to if(0))
		for i := 0; i < 10; i++ {
			// Execute immediately without spawning goroutine (undeferred)
			func() {
				variable++ // Sequential execution - no data race
			}()
		}
	}()

	wg.Wait()

	fmt.Printf("%d\n", variable)
}
