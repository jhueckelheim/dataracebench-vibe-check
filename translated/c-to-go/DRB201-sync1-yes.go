/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * This is a program based on a dataset contributed by
 * Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

 * Thread with id 1 acquires and releases the lock, but then it modifies x without holding it.
 * Data race pair: size@35:7:W vs. size@42:7:W
 */

package main

import (
	"fmt"
	"sync"
)

var (
	l sync.Mutex
	x = 1
)

func main() {
	var wg sync.WaitGroup
	barrier1 := make(chan bool)
	barrier2 := make(chan bool)

	wg.Add(2)

	// Thread 0
	go func() {
		defer wg.Done()

		// Wait for both threads to reach barrier
		barrier1 <- true
		<-barrier1

		l.Lock()
		x = 0
		l.Unlock()

		// Wait for both threads to reach second barrier
		barrier2 <- true
		<-barrier2
	}()

	// Thread 1
	go func() {
		defer wg.Done()

		// Wait for both threads to reach barrier
		<-barrier1
		barrier1 <- true

		l.Lock()
		l.Unlock()
		// Race condition: modifying x without holding the lock
		x = 1

		// Wait for both threads to reach second barrier
		<-barrier2
		barrier2 <- true
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
