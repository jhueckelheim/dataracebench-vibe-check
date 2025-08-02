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
	"sync"
)

// Module DRB102 translated to package-level variables
var y int
var x float64

func main() {
	//$omp parallel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//$omp single
		x = 1.0
		y = 1
		//$omp end single copyprivate(x,y)
	}()
	wg.Wait()
	//$omp end parallel

	fmt.Printf("x =%3.1f  y =%3d\n", x, y)
}
