//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//This example is modified version of nestable_lock.1.c example, OpenMP 5.0 Application Programming Examples.
//A nested lock can be locked several times. It doesn't unlock until you have unset it as many times as the
//number of calls to omp_set_nest_lock.
//incr_b is called at line 54 and line 59. So, it needs a nest_lock for p%b@35:5. No data race.

package main

import (
	"fmt"
	"sync"
)

type Pair struct {
	a   int
	b   int
	lck sync.Mutex
}

func incrA(p *Pair, a int) {
	p.a = p.a + 1
}

func incrB(p *Pair, b int) {
	p.lck.Lock()  // Nested lock protection
	p.b = p.b + 1 // No race - properly protected
	p.lck.Unlock()
}

func main() {
	var a, b int
	var p Pair

	p.a = 0
	p.b = 0

	//$omp parallel sections
	var wg sync.WaitGroup

	// Section 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.lck.Lock() // Outer lock
		incrB(&p, a) // This will acquire lock again (nested)
		incrA(&p, b)
		p.lck.Unlock() // Outer unlock
	}()

	// Section 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		incrB(&p, b) // This acquires its own lock
	}()

	wg.Wait()
	//$omp end parallel sections

	fmt.Printf("%d\n", p.b)
}
