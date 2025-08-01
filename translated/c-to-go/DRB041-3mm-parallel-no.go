/**
 * 3mm.c: This file is part of the PolyBench/C 3.2 test suite.
 * three steps of matrix multiplication to multiply four matrices.
 *
 * Contact: Louis-Noel Pouchet <pouchet@cse.ohio-state.edu>
 * Web address: http://polybench.sourceforge.net
 * License: /LICENSE.OSU.txt
 */

package main

import (
	"fmt"
	"os"
	"sync"
)

const NI = 128
const NJ = 128
const NK = 128
const NL = 128
const NM = 128

func initArray(ni, nj, nk, nl, nm int, A, B, C, D *[NI][NJ]float64) {
	var wg sync.WaitGroup
	for c1 := 0; c1 < ni; c1++ {
		wg.Add(1)
		go func(c1 int) {
			defer wg.Done()
			for c2 := 0; c2 < nj; c2++ {
				A[c1][c2] = float64(c1*c2) / float64(ni)
				B[c1][c2] = float64(c1*(c2+1)) / float64(nj)
				C[c1][c2] = float64(c1*(c2+3)) / float64(nl)
				D[c1][c2] = float64(c1*(c2+2)) / float64(nk)
			}
		}(c1)
	}
	wg.Wait()
}

func printArray(ni, nl int, G *[NI][NL]float64) {
	for i := 0; i < ni; i++ {
		for j := 0; j < nl; j++ {
			fmt.Printf("%.2f ", G[i][j])
			if (i*ni+j)%20 == 0 {
				fmt.Println()
			}
		}
	}
	fmt.Println()
}

func kernel3mm(ni, nj, nk, nl, nm int, E, A, B, F, C, D, G *[NI][NJ]float64) {
	// Initialize result matrices
	var wg1 sync.WaitGroup
	for c1 := 0; c1 < ni; c1++ {
		wg1.Add(1)
		go func(c1 int) {
			defer wg1.Done()
			for c2 := 0; c2 < nj; c2++ {
				G[c1][c2] = 0
				F[c1][c2] = 0
			}
		}(c1)
	}
	wg1.Wait()

	// F := C * D
	var wg2 sync.WaitGroup
	for c1 := 0; c1 < ni; c1++ {
		wg2.Add(1)
		go func(c1 int) {
			defer wg2.Done()
			for c2 := 0; c2 < nj; c2++ {
				for c5 := 0; c5 < nk; c5++ {
					F[c1][c2] += C[c1][c5] * D[c5][c2]
				}
			}
		}(c1)
	}
	wg2.Wait()

	// Initialize E
	var wg3 sync.WaitGroup
	for c1 := 0; c1 < ni; c1++ {
		wg3.Add(1)
		go func(c1 int) {
			defer wg3.Done()
			for c2 := 0; c2 < nj; c2++ {
				E[c1][c2] = 0
			}
		}(c1)
	}
	wg3.Wait()

	// E := A * B and G := E * F
	var wg4 sync.WaitGroup
	for c1 := 0; c1 < ni; c1++ {
		wg4.Add(1)
		go func(c1 int) {
			defer wg4.Done()
			for c2 := 0; c2 < nj; c2++ {
				// E[c1][c2] := A[c1][:] * B[:][c2]
				for c5 := 0; c5 < nk; c5++ {
					E[c1][c2] += A[c1][c5] * B[c5][c2]
				}
				// G[c1][:] += E[c1][c2] * F[c2][:]
				for c5 := 0; c5 < nl; c5++ {
					G[c1][c5] += E[c1][c2] * F[c2][c5]
				}
			}
		}(c1)
	}
	wg4.Wait()
}

func main() {
	ni := NI
	nj := NJ
	nk := NK
	nl := NL
	nm := NM

	var E [NI][NJ]float64
	var A [NI][NJ]float64
	var B [NI][NJ]float64
	var F [NI][NJ]float64
	var C [NI][NJ]float64
	var D [NI][NJ]float64
	var G [NI][NL]float64

	// Initialize arrays
	initArray(ni, nj, nk, nl, nm, &A, &B, &C, &D)

	// Run kernel
	kernel3mm(ni, nj, nk, nl, nm, &E, &A, &B, &F, &C, &D, &G)

	// Prevent dead-code elimination
	if len(os.Args) > 42 {
		printArray(ni, nl, &G)
	}
}
