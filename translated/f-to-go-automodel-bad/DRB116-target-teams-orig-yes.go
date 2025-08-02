//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//use of omp target + teams
//Without protection, master threads from two teams cause data races.
//Data race pair: a@24:9:W vs. a@24:9:W

package main

import (
	"fmt"
	"sync"
)

func main() {
	var i, length int
	var a []float64

	length = 100
	a = make([]float64, length)

	for i = 1; i <= length; i++ {
		a[i-1] = float64(i) / 2.0
	}

	//$omp target map(tofrom: a(0:len))
	//$omp teams num_teams(2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		a[49] = a[49] * 2.0
	}()
	go func() {
		defer wg.Done()
		a[49] = a[49] * 2.0
	}()
	wg.Wait()
	//$omp end teams
	//$omp end target

	fmt.Printf("a(50)= %f\n", a[49])
}
