/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
	This example is modified version of nestable_lock.1.c example, OpenMP 5.0 Application Programming Examples.

A nested lock can be locked several times. It doesn't unlock until you have unset
it as many times as the number of calls to lock.
incr_b is called from two places. So, it needs a nested lock for p.b.
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
)

type Pair struct {
	a   int
	b   int
	lck sync.RWMutex // Using RWMutex to simulate nested lock behavior
}

func (p *Pair) incr_a() {
	p.a++
}

func (p *Pair) incr_b() {
	p.lck.Lock()
	p.b++
	p.lck.Unlock()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	p := &Pair{a: 0, b: 0}

	var wg sync.WaitGroup
	wg.Add(2)

	// First section
	go func() {
		defer wg.Done()

		// Nested lock usage - lock then call function that also locks
		p.lck.Lock()
		p.incr_b() // This function also acquires the lock (nested behavior)
		p.incr_a()
		p.lck.Unlock()
	}()

	// Second section
	go func() {
		defer wg.Done()
		p.incr_b() // Direct call to protected function
	}()

	wg.Wait()

	fmt.Printf("%d\n", p.b)
}
