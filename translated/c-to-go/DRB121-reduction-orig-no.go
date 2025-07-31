/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
Number of threads is empirical: We need enough threads so that
the reduction is really performed hierarchically!
There is no data race.
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var variable int64 = 0
	numThreads := runtime.NumCPU()
	var wg sync.WaitGroup

	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Local reductions for each thread
			localSum1 := 0
			localSum2 := 0

			// First reduction: sum1 += i for i in 0..4
			for i := 0; i < 5; i++ {
				localSum1 += i
			}

			// Second reduction: sum2 += i for i in 0..4
			for i := 0; i < 5; i++ {
				localSum2 += i
			}

			// Combine local results
			localVar := localSum1 + localSum2 // Each thread contributes 10

			// Atomic reduction across threads
			atomic.AddInt64(&variable, int64(localVar))
		}(t)
	}

	wg.Wait()

	expectedResult := 20 * numThreads
	error := (int(variable) != expectedResult)

	if error {
		fmt.Printf("%d %d\n", variable, expectedResult)
		panic("Incorrect result!")
	}
}
