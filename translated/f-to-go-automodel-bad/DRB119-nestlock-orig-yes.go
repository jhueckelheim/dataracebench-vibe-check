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

// Module DRB118 translated to Go struct
type pair struct {
	a   int
	b   int
	lck sync.RWMutex
}

func incrA(p *pair, a int) {
	p.a = p.a + 1
}

func incrB(p *pair, b int) {
	p.b = p.b + 1
}

func main() {
	var a, b int

	var p pair
	p.a = 0
	p.b = 0

	//$omp parallel sections
	var wg sync.WaitGroup
	wg.Add(2)
	//$omp section
	go func() {
		defer wg.Done()
		p.lck.Lock()
		incrB(&p, a)
		incrA(&p, b)
		p.lck.Unlock()
	}()

	//$omp section
	go func() {
		defer wg.Done()
		incrB(&p, b)
	}()
	//$omp end parallel sections
	wg.Wait()

	fmt.Printf("%d\n", p.b)
}
