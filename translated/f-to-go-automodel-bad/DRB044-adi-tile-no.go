//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//*****************************************************************************
//
//  adi.F90: This file is part of the PolyBench/Fortran 1.0 test suite.
//
//  Contact: Louis-Noel Pouchet <pouchet@cse.ohio-state.edu>
//  Web address: http://polybench.sourceforge.net
//
//*****************************************************************************

//No data race pairs.

package main

import (
	"fmt"
	"sync"
)

type DATA_TYPE = float64

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func initArray(n int, x, a, b [][]DATA_TYPE) {
	if n >= 1 {
		//$omp parallel do private(c2,c3,c4)
		var wg sync.WaitGroup
		tileSize := 16
		numTiles := (n - 1) / tileSize
		if numTiles < 1 {
			numTiles = 1
		}

		for c1 := 1; c1 <= numTiles; c1++ {
			for c2 := 1; c2 <= numTiles; c2++ {
				wg.Add(1)
				go func(c1, c2 int) {
					defer wg.Done()
					for c3 := 16 * c1; c3 <= min(16*c1+15, n-1); c3++ {
						//$omp simd
						for c4 := 16 * c2; c4 <= min(16*c2+15, n-1); c4++ {
							x[c3-1][c4-1] = (float64(c3)*float64(c4+1) + 1.0) / float64(n)
							a[c3-1][c4-1] = (float64(c3)*float64(c4+2) + 2.0) / float64(n)
							b[c3-1][c4-1] = (float64(c3)*float64(c4+3) + 3.0) / float64(n)
						}
					}
				}(c1, c2)
			}
		}
		wg.Wait()
		//$omp end parallel do
	}
}

func printArray(n int, x [][]DATA_TYPE) {
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			fmt.Printf("%f ", x[i-1][j-1])
			if (i*500+j)%20 == 0 {
				fmt.Println()
			}
		}
	}
	fmt.Println()
}

func kernelAdi(tsteps, n int, x, a, b [][]DATA_TYPE) {
	if n >= 1 && tsteps >= 1 {
		tileSize := 16
		numTiles := (n - 1) / tileSize
		if numTiles < 1 {
			numTiles = 1
		}

		for c0 := 1; c0 <= tsteps; c0++ {
			if n >= 2 {
				//$omp parallel do private(c15,c9,c8)
				var wg1 sync.WaitGroup

				for c2 := 1; c2 <= numTiles; c2++ {
					wg1.Add(1)
					go func(c2 int) {
						defer wg1.Done()
						for c8 := 1; c8 <= numTiles; c8++ {
							for c9 := max(2, 16*c8); c9 <= min(16*c8+15, n-1); c9++ {
								//$omp simd
								for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
									b[c15-1][c9-1] = b[c15-1][c9-1] - a[c15-1][c9-1]*a[c15-1][c9-1]/b[c15-1][c9-2]
								}
							}
						}

						for c8 := 1; c8 <= numTiles; c8++ {
							for c9 := max(2, 16*c8); c9 <= min(16*c8+15, n-1); c9++ {
								//$omp simd
								for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
									x[c15-1][c9-1] = x[c15-1][c9-1] - x[c15-1][c9-1]*a[c15-1][c9-1]/b[c15-1][c9-2]
								}
							}
						}

						for c8 := 1; c8 <= (n-3)/tileSize; c8++ {
							for c9 := 16 * c8; c9 <= min(16*c8+15, n-3); c9++ {
								//$omp simd
								for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
									x[c15-1][n-c9-3] = (x[c15-1][n-3-c9] - x[c15-1][n-3-c9-1]*a[c15-1][n-c9-3]) / b[c15-1][n-3-c9]
								}
							}
						}
					}(c2)
				}
				wg1.Wait()
				//$omp end parallel do
			}

			//$omp parallel do private(c15)
			var wg2 sync.WaitGroup
			for c2 := 1; c2 <= numTiles; c2++ {
				wg2.Add(1)
				go func(c2 int) {
					defer wg2.Done()
					//$omp simd
					for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
						x[c15-1][n-2] = x[c15-1][n-2] / b[c15-1][n-2]
					}
				}(c2)
			}
			wg2.Wait()
			//$omp end parallel do

			if n >= 2 {
				//$omp parallel do private(c15, c9, c8)
				var wg3 sync.WaitGroup
				for c2 := 1; c2 <= numTiles; c2++ {
					wg3.Add(1)
					go func(c2 int) {
						defer wg3.Done()
						for c8 := 1; c8 <= numTiles; c8++ {
							for c9 := max(2, 16*c8); c9 <= min(16*c8+15, n-1); c9++ {
								//$omp simd
								for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
									b[c9-1][c15-1] = b[c9-1][c15-1] - a[c9-1][c15-1]*a[c9-1][c15-1]/b[c9-2][c15-1]
								}
							}
						}

						for c8 := 1; c8 <= numTiles; c8++ {
							for c9 := max(2, 16*c8); c9 <= min(16*c8+15, n-1); c9++ {
								//$omp simd
								for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
									x[c9-1][c15-1] = x[c9-1][c15-1] - x[c9-2][c15-1]*a[c9-1][c15-1]/b[c9-2][c15-1]
								}
							}
						}

						for c8 := 1; c8 <= (n-3)/tileSize; c8++ {
							for c9 := 16 * c8; c9 <= min(16*c8+15, n-3); c9++ {
								//$omp simd
								for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
									x[n-3-c9][c15-1] = (x[n-3-c9][c15-1] - x[n-c9-3][c15-1]*a[n-3-c9][c15-1]) / b[n-3-c9][c15-1]
								}
							}
						}
					}(c2)
				}
				wg3.Wait()
				//$omp end parallel do
			}

			//$omp parallel do private(c15)
			var wg4 sync.WaitGroup
			for c2 := 1; c2 <= numTiles; c2++ {
				wg4.Add(1)
				go func(c2 int) {
					defer wg4.Done()
					//$omp simd
					for c15 := 16 * c2; c15 <= min(16*c2+15, n-1); c15++ {
						x[n-2][c15-1] = x[n-2][c15-1] / b[n-2][c15-1]
					}
				}(c2)
			}
			wg4.Wait()
			//$omp end parallel do
		}
	}
}

func main() {
	var x, a, b [][]DATA_TYPE
	n := 500

	// Allocation of Arrays
	x = make([][]DATA_TYPE, n)
	a = make([][]DATA_TYPE, n)
	b = make([][]DATA_TYPE, n)
	for i := range x {
		x[i] = make([]DATA_TYPE, n)
		a[i] = make([]DATA_TYPE, n)
		b[i] = make([]DATA_TYPE, n)
	}

	// Initialization
	initArray(n, x, a, b)

	// Kernel Execution
	kernelAdi(10, n, x, a, b)

	// Prevent dead-code elimination. All live-out data must be printed
	// by the function call in argument.
	printArray(n, x)
}
