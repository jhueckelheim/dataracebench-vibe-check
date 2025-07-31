/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
A nested lock can be locked several times. It doesn't unlock until you have unset
it as many times as the number of calls to lock.
incr_b is called from two places. So, it needs a nested lock enclosing the function.
Missing nested lock will lead to race condition.
Data Race Pairs: p.b (write vs. write)
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
	lck sync.RWMutex
}

func (p *Pair) incr_a() {
	p.a++
}

func (p *Pair) incr_b() {
	// Missing lock protection - this causes the data race
	p.b++ // Data race: concurrent writes to p.b
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	p := &Pair{a: 0, b: 0}

	var wg sync.WaitGroup
	wg.Add(2)

	// First section
	go func() {
		defer wg.Done()

		p.lck.Lock()
		p.incr_b() // Data race: calls unprotected function while holding lock
		p.incr_a()
		p.lck.Unlock()
	}()

	// Second section
	go func() {
		defer wg.Done()
		p.incr_b() // Data race: calls unprotected function without lock
	}()

	wg.Wait()

	fmt.Printf("%d\n", p.b)
}
