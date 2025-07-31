/**
 * jacobi-2d-imper.c: This file is part of the PolyBench/C 3.2 test suite.
 * Jacobi with array copying, no reduction.
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

func initArray(n int, A, B *[N][N]float64) {
	var wg sync.WaitGroup
	for c1 := 0; c1 < n; c1++ {
		wg.Add(1)
		go func(c1 int) {
			defer wg.Done()
			for c2 := 0; c2 < n; c2++ {
				A[c1][c2] = (float64(c1)*(float64(c2)+2) + 2) / float64(n)
				B[c1][c2] = (float64(c1)*(float64(c2)+3) + 3) / float64(n)
			}
		}(c1)
	}
	wg.Wait()
}

func printArray(n int, A *[N][N]float64) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Fprintf(os.Stderr, "%0.2f ", A[i][j])
			if (i*n+j)%20 == 0 {
				fmt.Fprintf(os.Stderr, "\n")
			}
		}
	}
	fmt.Fprintf(os.Stderr, "\n")
}

func kernelJacobi2d(tsteps, n int, A, B *[N][N]float64) {
	// Simplified version of the complex loop structure
	for t := 0; t < tsteps; t++ {
		// First step: compute B from A in parallel
		var wg1 sync.WaitGroup
		for i := 1; i < n-1; i++ {
			wg1.Add(1)
			go func(i int) {
				defer wg1.Done()
				for j := 1; j < n-1; j++ {
					B[i][j] = 0.2 * (A[i][j] + A[i][j-1] + A[i][j+1] + A[i+1][j] + A[i-1][j])
				}
			}(i)
		}
		wg1.Wait()

		// Second step: copy B back to A in parallel
		var wg2 sync.WaitGroup
		for i := 1; i < n-1; i++ {
			wg2.Add(1)
			go func(i int) {
				defer wg2.Done()
				for j := 1; j < n-1; j++ {
					A[i][j] = B[i][j]
				}
			}(i)
		}
		wg2.Wait()
	}
}

func main() {
	n := N
	tsteps := TSTEPS

	var A [N][N]float64
	var B [N][N]float64

	// Initialize arrays
	initArray(n, &A, &B)

	// Run kernel
	kernelJacobi2d(tsteps, n, &A, &B)

	// Print result if requested
	if len(os.Args) > 42 && os.Args[0] == "" {
		printArray(n, &A)
	}
}
