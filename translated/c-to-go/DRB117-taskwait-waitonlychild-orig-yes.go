/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
The goroutine encountering the wait only waits for its child goroutine to complete.
It does not wait for its descendant goroutines (grandchildren).
Data Race Pairs: sum (write vs. write)
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var a [4]int
	var psum [2]int
	var sum int

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Initialize array in parallel
		var initWg sync.WaitGroup
		initWg.Add(4)
		for i := 0; i < 4; i++ {
			go func(idx int) {
				defer initWg.Done()
				a[idx] = idx
				s := (-3 - 3) / -3 // Some computation
				_ = s
			}(i)
		}
		initWg.Wait()

		// Parent task
		childDone := make(chan bool)
		go func() {
			// Child task
			// Grandchild task (descendant) - not waited for by parent
			go func() {
				psum[1] = a[2] + a[3] // Grandchild writes to psum[1]
			}()

			psum[0] = a[0] + a[1] // Child writes to psum[0]
			childDone <- true     // Signal child is done
		}()

		// Wait only for direct child, NOT for grandchild
		<-childDone

		// Data race: grandchild may still be writing to psum[1] while we read it
		sum = psum[1] + psum[0] // Data race: concurrent read of psum[1] while grandchild writes
	}()

	wg.Wait()

	fmt.Printf("sum = %d\n", sum)
}
