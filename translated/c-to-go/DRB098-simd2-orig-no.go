/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB098-simd2-orig-no.c

Description: Two-dimension array computation with a vectorization directive
collapse(2) makes simd associate with 2 loops.
Loop iteration variables should be predetermined as lastprivate.
*/

package main

import (
	"fmt"
)

func main() {
	const len = 100
	var a, b, c [len][len]float64

	// Initialize arrays
	for i := 0; i < len; i++ {
		for j := 0; j < len; j++ {
			a[i][j] = float64(i) / 2.0
			b[i][j] = float64(i) / 3.0
			c[i][j] = float64(i) / 7.0
		}
	}

	// Simulate SIMD collapse(2) - Go compiler will auto-vectorize when possible
	// The collapse(2) flattens both loops for SIMD processing
	// In Go, we trust the compiler's auto-vectorization capabilities
	for i := 0; i < len; i++ {
		for j := 0; j < len; j++ {
			c[i][j] = a[i][j] * b[i][j]
		}
	}

	fmt.Printf("c[50][50]=%f\n", c[50][50])
}
