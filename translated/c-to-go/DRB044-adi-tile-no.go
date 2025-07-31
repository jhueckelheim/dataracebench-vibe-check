/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * adi.c: This file is part of the PolyBench/C 3.2 test suite.
 * Alternating Direction Implicit solver with tiling and nested SIMD.
 * Race-free version with proper synchronization.
 */

package main

import (
	"fmt"
	"sync"
)

const (
	N        = 128 // Reduced size for practical Go execution
	TSTEPS   = 10
	TILESIZE = 16
)

func initArray(n int, X, A, B [][]float64) {
	tilesPerRow := (n + TILESIZE - 1) / TILESIZE
	var wg sync.WaitGroup

	for tileI := 0; tileI < tilesPerRow; tileI++ {
		for tileJ := 0; tileJ < tilesPerRow; tileJ++ {
			wg.Add(1)
			go func(ti, tj int) {
				defer wg.Done()

				iStart := ti * TILESIZE
				iEnd := iStart + TILESIZE
				if iEnd > n {
					iEnd = n
				}

				jStart := tj * TILESIZE
				jEnd := jStart + TILESIZE
				if jEnd > n {
					jEnd = n
				}

				for i := iStart; i < iEnd; i++ {
					for j := jStart; j < jEnd; j++ {
						X[i][j] = (float64(i)*float64(j+1) + 1) / float64(n)
						A[i][j] = (float64(i)*float64(j+2) + 2) / float64(n)
						B[i][j] = (float64(i)*float64(j+3) + 3) / float64(n)
					}
				}
			}(tileI, tileJ)
		}
	}
	wg.Wait()
}

func kernelADI(tsteps, n int, X, A, B [][]float64) {
	tilesPerRow := (n + TILESIZE - 1) / TILESIZE

	for t := 0; t < tsteps; t++ {
		if n >= 2 {
			// First pass: row-wise operations
			var wg1 sync.WaitGroup
			for tileI := 0; tileI < tilesPerRow; tileI++ {
				wg1.Add(1)
				go func(ti int) {
					defer wg1.Done()

					iStart := ti * TILESIZE
					iEnd := iStart + TILESIZE
					if iEnd > n {
						iEnd = n
					}

					for tileJ := 0; tileJ < tilesPerRow; tileJ++ {
						jStart := tileJ * TILESIZE
						jEnd := jStart + TILESIZE
						if jEnd > n {
							jEnd = n
						}

						// B update - row operations
						for j := jStart; j < jEnd; j++ {
							if j > 0 {
								for i := iStart; i < iEnd; i++ {
									B[i][j] = B[i][j] - A[i][j]*A[i][j]/B[i][j-1]
								}
							}
						}

						// X update - row operations
						for j := jStart; j < jEnd; j++ {
							if j > 0 {
								for i := iStart; i < iEnd; i++ {
									X[i][j] = X[i][j] - X[i][j-1]*A[i][j]/B[i][j-1]
								}
							}
						}

						// Backward elimination - row operations
						for j := jEnd - 1; j >= jStart; j-- {
							if j < n-2 {
								for i := iStart; i < iEnd; i++ {
									X[i][j] = (X[i][j] - X[i][j+1]*A[i][j]) / B[i][j]
								}
							}
						}
					}
				}(tileI)
			}
			wg1.Wait()
		}

		// Final row operation
		var wgFinal1 sync.WaitGroup
		for tileI := 0; tileI < tilesPerRow; tileI++ {
			wgFinal1.Add(1)
			go func(ti int) {
				defer wgFinal1.Done()

				iStart := ti * TILESIZE
				iEnd := iStart + TILESIZE
				if iEnd > n {
					iEnd = n
				}

				for i := iStart; i < iEnd; i++ {
					X[i][n-1] = X[i][n-1] / B[i][n-1]
				}
			}(tileI)
		}
		wgFinal1.Wait()

		if n >= 2 {
			// Second pass: column-wise operations
			var wg2 sync.WaitGroup
			for tileJ := 0; tileJ < tilesPerRow; tileJ++ {
				wg2.Add(1)
				go func(tj int) {
					defer wg2.Done()

					jStart := tj * TILESIZE
					jEnd := jStart + TILESIZE
					if jEnd > n {
						jEnd = n
					}

					for tileI := 0; tileI < tilesPerRow; tileI++ {
						iStart := tileI * TILESIZE
						iEnd := iStart + TILESIZE
						if iEnd > n {
							iEnd = n
						}

						// B update - column operations
						for i := iStart; i < iEnd; i++ {
							if i > 0 {
								for j := jStart; j < jEnd; j++ {
									B[i][j] = B[i][j] - A[i][j]*A[i][j]/B[i-1][j]
								}
							}
						}

						// X update - column operations
						for i := iStart; i < iEnd; i++ {
							if i > 0 {
								for j := jStart; j < jEnd; j++ {
									X[i][j] = X[i][j] - X[i-1][j]*A[i][j]/B[i-1][j]
								}
							}
						}

						// Backward elimination - column operations
						for i := iEnd - 1; i >= iStart; i-- {
							if i < n-2 {
								for j := jStart; j < jEnd; j++ {
									X[i][j] = (X[i][j] - X[i+1][j]*A[i][j]) / B[i][j]
								}
							}
						}
					}
				}(tileJ)
			}
			wg2.Wait()
		}

		// Final column operation
		var wgFinal2 sync.WaitGroup
		for tileJ := 0; tileJ < tilesPerRow; tileJ++ {
			wgFinal2.Add(1)
			go func(tj int) {
				defer wgFinal2.Done()

				jStart := tj * TILESIZE
				jEnd := jStart + TILESIZE
				if jEnd > n {
					jEnd = n
				}

				for j := jStart; j < jEnd; j++ {
					X[n-1][j] = X[n-1][j] / B[n-1][j]
				}
			}(tileJ)
		}
		wgFinal2.Wait()
	}
}

func printPartialArray(n int, X [][]float64) {
	fmt.Printf("Sample values: X[0][0]=%.6f, X[%d][%d]=%.6f, X[%d][%d]=%.6f\n",
		X[0][0], n/2, n/2, X[n/2][n/2], n-1, n-1, X[n-1][n-1])
}

func main() {
	n := N
	tsteps := TSTEPS

	// Allocate arrays
	X := make([][]float64, n)
	A := make([][]float64, n)
	B := make([][]float64, n)
	for i := range X {
		X[i] = make([]float64, n)
		A[i] = make([]float64, n)
		B[i] = make([]float64, n)
	}

	// Initialize arrays
	initArray(n, X, A, B)

	// Run ADI kernel
	kernelADI(tsteps, n, X, A, B)

	// Print results
	printPartialArray(n, X)
}
