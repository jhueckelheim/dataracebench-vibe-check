//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//****************************************************************************
//
//  adi.F90: This file is part of the PolyBench/Fortran 1.0 test suite.
//
//  Contact: Louis-Noel Pouchet <pouchet@cse.ohio-state.edu>
//  Web address: http://polybench.sourceforge.net
//
//****************************************************************************

//No data race pairs.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

const N = 500

func main() {
	// Array declarations (replacing polybench macros)
	var x, a, b [N][N]float64

	// Initialization
	initArray(N, &x, &a, &b)

	// Kernel Execution
	kernelAdi(10, N, &x, &a, &b)

	// Print result to prevent dead-code elimination
	printArray(N, &x)
}

func initArray(n int, x, a, b *[N][N]float64) {
	if n >= 1 {
		//$omp parallel do private(c2)
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := (n - 1) / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 1; start <= n-1; start += chunkSize {
			end := start + chunkSize - 1
			if end > n-1 {
				end = n - 1
			}
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for c1 := start; c1 <= end; c1++ {
					for c2 := 1; c2 <= n-1; c2++ {
						x[c1-1][c2-1] = (float64(c1)*float64(c2+1) + 1.0) / float64(n)
						a[c1-1][c2-1] = (float64(c1)*float64(c2+2) + 2.0) / float64(n)
						b[c1-1][c2-1] = (float64(c1)*float64(c2+3) + 3.0) / float64(n)
					}
				}
			}(start, end)
		}
		wg.Wait()
		//$omp end parallel do
	}
}

func printArray(n int, x *[N][N]float64) {
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			fmt.Printf("%0.6f ", x[i-1][j-1])
			if (i*500+j)%20 == 0 {
				fmt.Println()
			}
		}
	}
	fmt.Println()
}

func kernelAdi(tsteps, n int, x, a, b *[N][N]float64) {
	for c0 := 1; c0 <= 10; c0++ {
		//$omp parallel do private(c8)
		var wg1 sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := 500 / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 1; start <= 500; start += chunkSize {
			end := start + chunkSize - 1
			if end > 500 {
				end = 500
			}
			wg1.Add(1)
			go func(start, end int) {
				defer wg1.Done()
				for c2 := start; c2 <= end; c2++ {
					for c8 := 2; c8 <= 500; c8++ {
						b[c2-1][c8-1] = b[c2-1][c8-1] - a[c2-1][c8-1]*a[c2-1][c8-1]/b[c2-1][c8-2]
					}

					for c8 := 2; c8 <= 500; c8++ {
						x[c2-1][c8-1] = x[c2-1][c8-1] - x[c2-1][c8-2]*a[c2-1][c8-1]/b[c2-1][c8-2]
					}

					for c8 := 1; c8 <= 498; c8++ {
						x[c2-1][500-c8-1] = (x[c2-1][500-c8-1] - x[c2-1][500-c8-2]*a[c2-1][500-c8-2]) / b[c2-1][500-2-c8]
					}
				}
			}(start, end)
		}
		wg1.Wait()
		//$omp end parallel do

		//$omp parallel do
		var wg2 sync.WaitGroup
		for start := 1; start <= 500; start += chunkSize {
			end := start + chunkSize - 1
			if end > 500 {
				end = 500
			}
			wg2.Add(1)
			go func(start, end int) {
				defer wg2.Done()
				for c2 := start; c2 <= end; c2++ {
					x[c2-1][498] = x[c2-1][498] / b[c2-1][498]
				}
			}(start, end)
		}
		wg2.Wait()
		//$omp end parallel do

		//$omp parallel do private(c8)
		var wg3 sync.WaitGroup
		for start := 1; start <= 500; start += chunkSize {
			end := start + chunkSize - 1
			if end > 500 {
				end = 500
			}
			wg3.Add(1)
			go func(start, end int) {
				defer wg3.Done()
				for c2 := start; c2 <= end; c2++ {
					for c8 := 2; c8 <= 500; c8++ {
						b[c8-1][c2-1] = b[c8-1][c2-1] - a[c8-1][c2-1]*a[c8-1][c2-1]/b[c8-2][c2-1]
					}

					for c8 := 2; c8 <= 500; c8++ {
						x[c8-1][c2-1] = x[c8-1][c2-1] - x[c8-2][c2-1]*a[c8-1][c2-1]/b[c8-2][c2-1]
					}

					for c8 := 1; c8 <= 498; c8++ {
						x[500-c8-1][c2-1] = (x[500-c8-1][c2-1] - x[500-c8-2][c2-1]*a[500-2-c8][c2-1]) / b[500-c8-1][c2-1]
					}
				}
			}(start, end)
		}
		wg3.Wait()
		//$omp end parallel do

		//$omp parallel do
		var wg4 sync.WaitGroup
		for start := 1; start <= 500; start += chunkSize {
			end := start + chunkSize - 1
			if end > 500 {
				end = 500
			}
			wg4.Add(1)
			go func(start, end int) {
				defer wg4.Done()
				for c2 := start; c2 <= end; c2++ {
					x[498][c2-1] = x[498][c2-1] / b[498][c2-1]
				}
			}(start, end)
		}
		wg4.Wait()
		//$omp end parallel do
	}
}