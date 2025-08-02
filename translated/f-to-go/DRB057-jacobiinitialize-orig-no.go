//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Use of private() clause. No data race pairs.

package main

import (
	"runtime"
	"sync"
)

// Package-level variables (replacing module)
var MSIZE int
var n, m, mits int
var u, f, uold [][]float64
var dx, dy, tol, relax, alpha float64

func initialize() {
	MSIZE = 200
	mits = 1000
	relax = 1.0
	alpha = 0.0543
	n = MSIZE
	m = MSIZE

	// Allocate arrays
	u = make([][]float64, MSIZE)
	f = make([][]float64, MSIZE)
	uold = make([][]float64, MSIZE)
	for i := 0; i < MSIZE; i++ {
		u[i] = make([]float64, MSIZE)
		f[i] = make([]float64, MSIZE)
		uold[i] = make([]float64, MSIZE)
	}

	dx = 2.0 / float64(n-1)
	dy = 2.0 / float64(m-1)

	// Initialize initial condition and RHS
	//$omp parallel do private(i,j,xx,yy)
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := n / numCPU
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 1; start <= n; start += chunkSize {
		end := start + chunkSize - 1
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i <= end; i++ {
				for j := 1; j <= m; j++ {
					// All variables are private to each goroutine
					xx := int(-1.0 + dx*float64(i-1))
					yy := int(-1.0 + dy*float64(i-1))
					u[i-1][j-1] = 0.0
					xxf := float64(xx)
					yyf := float64(yy)
					f[i-1][j-1] = -1.0*alpha*(1.0-xxf*xxf)*(1.0-yyf*yyf) - 2.0*(1.0-xxf*xxf) - 2.0*(1.0-yyf*yyf)
				}
			}
		}(start, end)
	}
	wg.Wait()
	//$omp end parallel do
}

func main() {
	initialize()
}
