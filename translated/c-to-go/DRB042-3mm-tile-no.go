/*
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
!!! Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
!!! and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
!!!
!!! SPDX-License-Identifier: (BSD-3-Clause)
!!!~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!!!
*/

/*
 * 3mm.c: This file is part of the PolyBench/C 3.2 test suite.
 * with tiling 16x16 and nested SIMD
 * Race-free version with proper synchronization.
 */

package main

import (
	"fmt"
	"sync"
)

const (
	NI       = 128 // Reduced size for practical Go execution
	NJ       = 128
	NK       = 128
	NL       = 128
	NM       = 128
	TILESIZE = 16
)

func initArray(ni, nj, nk, nl, nm int, A, B, C, D [][]float64) {
	tilesPerRow := (ni + TILESIZE - 1) / TILESIZE
	tilesPerCol := (nj + TILESIZE - 1) / TILESIZE
	var wg sync.WaitGroup

	for tileI := 0; tileI < tilesPerRow; tileI++ {
		for tileJ := 0; tileJ < tilesPerCol; tileJ++ {
			wg.Add(1)
			go func(ti, tj int) {
				defer wg.Done()

				iStart := ti * TILESIZE
				iEnd := iStart + TILESIZE
				if iEnd > ni {
					iEnd = ni
				}

				jStart := tj * TILESIZE
				jEnd := jStart + TILESIZE
				if jEnd > nj {
					jEnd = nj
				}

				for i := iStart; i < iEnd; i++ {
					for j := jStart; j < jEnd; j++ {
						if j < len(A[i]) {
							A[i][j] = float64(i*j) / float64(ni)
						}
						if j < len(B[i]) {
							B[i][j] = float64(i*(j+1)) / float64(nj)
						}
						if j < len(C[i]) {
							C[i][j] = float64(i*(j+3)) / float64(nl)
						}
						if j < len(D[i]) {
							D[i][j] = float64(i*(j+2)) / float64(nk)
						}
					}
				}
			}(tileI, tileJ)
		}
	}
	wg.Wait()
}

func kernel3mm(ni, nj, nk, nl, nm int, E, A, B, F, C, D, G [][]float64) {
	// E := A * B
	computeMatrixMult(ni, nj, nk, E, A, B)

	// F := C * D
	computeMatrixMult(nj, nl, nm, F, C, D)

	// G := E * F
	computeMatrixMult(ni, nl, nj, G, E, F)
}

func computeMatrixMult(rows, cols, common int, result, matA, matB [][]float64) {
	tilesPerRow := (rows + TILESIZE - 1) / TILESIZE
	tilesPerCol := (cols + TILESIZE - 1) / TILESIZE
	var wg sync.WaitGroup

	for tileI := 0; tileI < tilesPerRow; tileI++ {
		for tileJ := 0; tileJ < tilesPerCol; tileJ++ {
			wg.Add(1)
			go func(ti, tj int) {
				defer wg.Done()

				iStart := ti * TILESIZE
				iEnd := iStart + TILESIZE
				if iEnd > rows {
					iEnd = rows
				}

				jStart := tj * TILESIZE
				jEnd := jStart + TILESIZE
				if jEnd > cols {
					jEnd = cols
				}

				// Initialize result tile to zero
				for i := iStart; i < iEnd; i++ {
					for j := jStart; j < jEnd; j++ {
						if i < len(result) && j < len(result[i]) {
							result[i][j] = 0.0
						}
					}
				}

				// Compute matrix multiplication for this tile
				for kTile := 0; kTile < (common+TILESIZE-1)/TILESIZE; kTile++ {
					kStart := kTile * TILESIZE
					kEnd := kStart + TILESIZE
					if kEnd > common {
						kEnd = common
					}

					for i := iStart; i < iEnd; i++ {
						for j := jStart; j < jEnd; j++ {
							if i < len(result) && j < len(result[i]) {
								for k := kStart; k < kEnd; k++ {
									if k < len(matA[i]) && k < len(matB) && j < len(matB[k]) {
										result[i][j] += matA[i][k] * matB[k][j]
									}
								}
							}
						}
					}
				}
			}(tileI, tileJ)
		}
	}
	wg.Wait()
}

func printPartialArray(ni, nl int, G [][]float64) {
	if ni > 0 && nl > 0 && len(G) > 0 && len(G[0]) > 0 {
		fmt.Printf("Sample values: G[0][0]=%.6f", G[0][0])
		if ni/2 < len(G) && nl/2 < len(G[ni/2]) {
			fmt.Printf(", G[%d][%d]=%.6f", ni/2, nl/2, G[ni/2][nl/2])
		}
		if ni-1 < len(G) && nl-1 < len(G[ni-1]) {
			fmt.Printf(", G[%d][%d]=%.6f", ni-1, nl-1, G[ni-1][nl-1])
		}
		fmt.Println()
	}
}

func main() {
	ni, nj, nk, nl, nm := NI, NJ, NK, NL, NM

	// Allocate arrays
	A := make([][]float64, ni)
	B := make([][]float64, nk)
	C := make([][]float64, nj)
	D := make([][]float64, nm)
	E := make([][]float64, ni) // A * B result
	F := make([][]float64, nj) // C * D result
	G := make([][]float64, ni) // E * F result

	for i := range A {
		A[i] = make([]float64, nk)
	}
	for i := range B {
		B[i] = make([]float64, nj)
	}
	for i := range C {
		C[i] = make([]float64, nm)
	}
	for i := range D {
		D[i] = make([]float64, nl)
	}
	for i := range E {
		E[i] = make([]float64, nj)
	}
	for i := range F {
		F[i] = make([]float64, nl)
	}
	for i := range G {
		G[i] = make([]float64, nl)
	}

	// Initialize arrays
	initArray(ni, nj, nk, nl, nm, A, B, C, D)

	// Run 3mm kernel
	kernel3mm(ni, nj, nk, nl, nm, E, A, B, F, C, D, G)

	// Print results
	printPartialArray(ni, nl, G)
}
