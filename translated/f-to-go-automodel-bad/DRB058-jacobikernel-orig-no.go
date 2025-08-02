//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Two parallel for loops within one single parallel region,
//combined with private() and reduction().

//3.7969326424804763E-007 vs 3.7969326424804758E-007. There is no race condition. The minute
//difference at 22nd point after decimal is due to the precision in fortran95

package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

// Module DRB058 translated to package-level variables
var MSIZE int
var n, m, mits int
var u, f, uold [][]float64
var dx, dy, tol, relax, alpha float64

func initialize() {
	var xx, yy int

	MSIZE = 200
	mits = 1000
	tol = 0.0000000001
	relax = 1.0
	alpha = 0.0543
	n = MSIZE
	m = MSIZE
	u = make([][]float64, MSIZE)
	f = make([][]float64, MSIZE)
	uold = make([][]float64, MSIZE)
	for i := range u {
		u[i] = make([]float64, MSIZE)
		f[i] = make([]float64, MSIZE)
		uold[i] = make([]float64, MSIZE)
	}

	dx = 2.0 / float64(n-1)
	dy = 2.0 / float64(m-1)

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			xx = int(-1.0 + dx*float64(i-1))
			yy = int(-1.0 + dy*float64(i-1))
			u[i-1][j-1] = 0.0
			f[i-1][j-1] = -1.0*alpha*(1.0-float64(xx*xx))*(1.0-float64(yy*yy)) - 2.0*(1.0-float64(xx*xx)) - 2.0*(1.0-float64(yy*yy))
		}
	}
}

func jacobi() {
	var omega float64
	var k int
	var error, resid, ax, ay, b float64

	MSIZE = 200
	mits = 1000
	tol = 0.0000000001
	relax = 1.0
	alpha = 0.0543
	n = MSIZE
	m = MSIZE

	omega = relax
	dx = 2.0 / float64(n-1)
	dy = 2.0 / float64(m-1)

	ax = 1.0 / (dx * dx) // X-direction coef
	ay = 1.0 / (dy * dy) // Y-direction coef
	b = -2.0/(dx*dx) - 2.0/(dy*dy) - alpha

	error = 10.0 * tol
	k = 1

	for k = 1; k <= mits; k++ {
		error = 0.0

		//Copy new solution into old
		//$omp parallel
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()

		//$omp do private(i,j)
		chunkSize1 := n / numCPU
		if chunkSize1 < 1 {
			chunkSize1 = 1
		}
		for start := 1; start <= n; start += chunkSize1 {
			end := start + chunkSize1 - 1
			if end > n {
				end = n
			}
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for i := start; i <= end; i++ {
					for j := 1; j <= m; j++ {
						uold[i-1][j-1] = u[i-1][j-1]
					}
				}
			}(start, end)
		}
		wg.Wait()
		//$omp end do

		//$omp do private(i,j,resid) reduction(+:error)
		var errorMutex sync.Mutex
		chunkSize2 := (n - 2) / numCPU
		if chunkSize2 < 1 {
			chunkSize2 = 1
		}
		for start := 2; start <= n-1; start += chunkSize2 {
			end := start + chunkSize2 - 1
			if end > n-1 {
				end = n - 1
			}
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				localError := 0.0
				for i := start; i <= end; i++ {
					for j := 2; j <= m-1; j++ {
						resid = (ax*(uold[i-2][j-1]+uold[i][j-1]) + ay*(uold[i-1][j-2]+uold[i-1][j]) + b*uold[i-1][j-1] - f[i-1][j-1]) / b
						u[i-1][j-1] = uold[i-1][j-1] - omega*resid
						localError += resid * resid
					}
				}
				errorMutex.Lock()
				error += localError
				errorMutex.Unlock()
			}(start, end)
		}
		wg.Wait()
		//$omp end do nowait
		//$omp end parallel

		//Error check
		error = math.Sqrt(error) / float64(n*m)
	}

	fmt.Printf("Total number of iterations: %d\n", k)
	fmt.Printf("Residual: %f\n", error)
}

func main() {
	initialize()
	jacobi()
}
