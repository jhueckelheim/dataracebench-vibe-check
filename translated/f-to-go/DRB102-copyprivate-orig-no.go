//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//threadprivate+copyprivate: no data races

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	// Simulate copyprivate: values set in single are copied to all threads
	var sharedX float64
	var sharedY int

	for threadID := 0; threadID < numCPU; threadID++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			// Each thread has its own copy (threadprivate equivalent)
			var x float64
			var y int

			//$omp single
			if threadID == 0 { // Only master thread executes single
				x = 1.0
				y = 1
				sharedX = x // Store for copyprivate
				sharedY = y
			}
			//$omp end single copyprivate(x,y)

			// Simulate copyprivate: all threads get copies of the values
			x = sharedX
			y = sharedY

		}(threadID)
	}
	wg.Wait()
	//$omp end parallel

	// Values from the copyprivate operation
	fmt.Printf("x = %3.1f  y = %3d\n", sharedX, sharedY)
}
