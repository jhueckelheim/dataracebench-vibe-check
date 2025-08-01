/**
 * adi.c: This file is part of the PolyBench/C 3.2 test suite.
 *
 * Alternating Direction Implicit solver:
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

const N = 500
const TSTEPS = 10

func initArray(n int, X *[N][N]float64, A *[N][N]float64, B *[N][N]float64) {
	var wg sync.WaitGroup
	for c1 := 0; c1 < n; c1++ {
		wg.Add(1)
		go func(c1 int) {
			defer wg.Done()
			for c2 := 0; c2 < n; c2++ {
				X[c1][c2] = (float64(c1)*(float64(c2)+1) + 1) / float64(n)
				A[c1][c2] = (float64(c1)*(float64(c2)+2) + 2) / float64(n)
				B[c1][c2] = (float64(c1)*(float64(c2)+3) + 3) / float64(n)
			}
		}(c1)
	}
	wg.Wait()
}

func printArray(n int, X *[N][N]float64) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%.2f ", X[i][j])
			if (i*N+j)%20 == 0 {
				fmt.Println()
			}
		}
	}
	fmt.Println()
}

func kernelAdi(tsteps int, n int, X *[N][N]float64, A *[N][N]float64, B *[N][N]float64) {
	for c0 := 0; c0 < tsteps; c0++ {
		var wg1 sync.WaitGroup
		for c2 := 0; c2 < n; c2++ {
			wg1.Add(1)
			go func(c2 int) {
				defer wg1.Done()
				for c8 := 1; c8 < n; c8++ {
					B[c2][c8] = B[c2][c8] - A[c2][c8]*A[c2][c8]/B[c2][c8-1]
				}
				for c8 := 1; c8 < n; c8++ {
					X[c2][c8] = X[c2][c8] - X[c2][c8-1]*A[c2][c8]/B[c2][c8-1]
				}
				for c8 := 0; c8 < n-2; c8++ {
					X[c2][n-c8-2] = (X[c2][n-2-c8] - X[c2][n-2-c8-1]*A[c2][n-c8-3]) / B[c2][n-3-c8]
				}
			}(c2)
		}
		wg1.Wait()

		var wg2 sync.WaitGroup
		for c2 := 0; c2 < n; c2++ {
			wg2.Add(1)
			go func(c2 int) {
				defer wg2.Done()
				X[c2][n-1] = X[c2][n-1] / B[c2][n-1]
			}(c2)
		}
		wg2.Wait()

		var wg3 sync.WaitGroup
		for c2 := 0; c2 < n; c2++ {
			wg3.Add(1)
			go func(c2 int) {
				defer wg3.Done()
				for c8 := 1; c8 < n; c8++ {
					B[c8][c2] = B[c8][c2] - A[c8][c2]*A[c8][c2]/B[c8-1][c2]
				}
				for c8 := 1; c8 < n; c8++ {
					X[c8][c2] = X[c8][c2] - X[c8-1][c2]*A[c8][c2]/B[c8-1][c2]
				}
				for c8 := 0; c8 < n-2; c8++ {
					X[n-2-c8][c2] = (X[n-2-c8][c2] - X[n-c8-3][c2]*A[n-3-c8][c2]) / B[n-2-c8][c2]
				}
			}(c2)
		}
		wg3.Wait()

		var wg4 sync.WaitGroup
		for c2 := 0; c2 < n; c2++ {
			wg4.Add(1)
			go func(c2 int) {
				defer wg4.Done()
				X[n-1][c2] = X[n-1][c2] / B[n-1][c2]
			}(c2)
		}
		wg4.Wait()
	}
}

func main() {
	n := N
	tsteps := TSTEPS

	var X [N][N]float64
	var A [N][N]float64
	var B [N][N]float64

	// Initialize arrays
	initArray(n, &X, &A, &B)

	// Run kernel
	kernelAdi(tsteps, n, &X, &A, &B)

	// Prevent dead-code elimination with conditional print
	if len(os.Args) > 42 {
		printArray(n, &X)
	}
}
