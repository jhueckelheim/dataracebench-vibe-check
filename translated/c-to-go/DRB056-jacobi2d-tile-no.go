/**
 * jacobi-2d-imper.c: This file is part of the PolyBench/C 3.2 test suite.
 * Jacobi with array copying, no reduction. with tiling and nested parallelization.
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
const TILE_SIZE = 16

func initArray(n int, A, B *[N][N]float64) {
	var wg sync.WaitGroup
	// Tiled initialization with nested parallelization
	for c1 := 0; c1 < (n+TILE_SIZE-1)/TILE_SIZE; c1++ {
		wg.Add(1)
		go func(c1 int) {
			defer wg.Done()
			for c2 := 0; c2 < (n+TILE_SIZE-1)/TILE_SIZE; c2++ {
				for c3 := c2 * TILE_SIZE; c3 < min((c2+1)*TILE_SIZE, n); c3++ {
					for c4 := c1 * TILE_SIZE; c4 < min((c1+1)*TILE_SIZE, n); c4++ {
						A[c4][c3] = (float64(c4)*(float64(c3)+2) + 2) / float64(n)
						B[c4][c3] = (float64(c4)*(float64(c3)+3) + 3) / float64(n)
					}
				}
			}
		}(c1)
	}
	wg.Wait()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

func kernelJacobi2dTiled(tsteps, n int, A, B *[N][N]float64) {
	// Simplified tiled version
	for t := 0; t < tsteps; t++ {
		// Parallel tiled computation
		var wg sync.WaitGroup
		for ti := 0; ti < (n+TILE_SIZE-1)/TILE_SIZE; ti++ {
			wg.Add(1)
			go func(ti int) {
				defer wg.Done()
				for tj := 0; tj < (n+TILE_SIZE-1)/TILE_SIZE; tj++ {
					// Process tile
					for i := max(ti*TILE_SIZE, 1); i < min((ti+1)*TILE_SIZE, n-1); i++ {
						for j := max(tj*TILE_SIZE, 1); j < min((tj+1)*TILE_SIZE, n-1); j++ {
							B[i][j] = 0.2 * (A[i][j] + A[i][j-1] + A[i][j+1] + A[i+1][j] + A[i-1][j])
						}
					}
				}
			}(ti)
		}
		wg.Wait()

		// Copy B back to A in parallel tiles
		var wg2 sync.WaitGroup
		for ti := 0; ti < (n+TILE_SIZE-1)/TILE_SIZE; ti++ {
			wg2.Add(1)
			go func(ti int) {
				defer wg2.Done()
				for tj := 0; tj < (n+TILE_SIZE-1)/TILE_SIZE; tj++ {
					for i := max(ti*TILE_SIZE, 1); i < min((ti+1)*TILE_SIZE, n-1); i++ {
						for j := max(tj*TILE_SIZE, 1); j < min((tj+1)*TILE_SIZE, n-1); j++ {
							A[i][j] = B[i][j]
						}
					}
				}
			}(ti)
		}
		wg2.Wait()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	n := N
	tsteps := TSTEPS

	var A [N][N]float64
	var B [N][N]float64

	// Initialize arrays
	initArray(n, &A, &B)

	// Run kernel
	kernelJacobi2dTiled(tsteps, n, &A, &B)

	// Print result if requested
	if len(os.Args) > 42 && os.Args[0] == "" {
		printArray(n, &A)
	}
}
