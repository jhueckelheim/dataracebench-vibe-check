/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

No race. The array b is divided into two non-overlapping halves that are referenced by u[0] and u[1].
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var (
	c      = 0.2
	n      = 20
	nsteps = 100
)

func main() {
	b := make([]float64, 2*n)
	u := [2][]float64{b[0:n], b[n : 2*n]} // Non-overlapping halves

	// Initialize arrays
	for i := 1; i < n-1; i++ {
		val := rand.Float64()
		u[0][i] = val
		u[1][i] = val
	}
	u[0][0] = 0.5
	u[0][n-1] = 0.5
	u[1][0] = 0.5
	u[1][n-1] = 0.5

	p := 0
	for t := 0; t < nsteps; t++ {
		var wg sync.WaitGroup

		// Parallel for loop
		for i := 1; i < n-1; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				// No race: u[0] and u[1] point to non-overlapping array regions
				u[1-p][idx] = u[p][idx] + c*(u[p][idx-1]+u[p][idx+1]-2*u[p][idx])
			}(i)
		}

		wg.Wait()
		p = 1 - p
	}

	// Print results
	for i := 0; i < n; i++ {
		fmt.Printf("%.2f ", u[p][i])
	}
	fmt.Println()
}
