//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//A nested lock can be locked several times. It doesn't unlock until you have unset
//it as many times as the number of calls to omp_set_nest_lock.
//incr_b is called at line 53 and line 58. So, it needs a nest_lock enclosing line 35
//Missing nest_lock will lead to race condition at line:35.
//Data Race Pairs, p%b@35:5:W vs. p%b@35:5:W

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
	// MISSING: Lock protection!
	p.b = p.b + 1 // RACE: Multiple sections access p.b without protection
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
		p.lck.Lock() // This lock doesn't protect incr_b
		incrB(&p, a) // RACE: incr_b has no lock protection
		incrA(&p, b)
		p.lck.Unlock()
	}()

	// Section 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		incrB(&p, b) // RACE: incr_b has no lock protection
	}()

	wg.Wait()
	//$omp end parallel sections

	fmt.Printf("%d\n", p.b)
}
