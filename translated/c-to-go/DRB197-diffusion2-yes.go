/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

Overlap of the two ranges u[0] and u[1] when u[1][i] is accessed.
Data race pairs: u[1 - p][i]@38:7:W vs. u[p][i - 1]@38:15:R
                 u[1 - p][i]@38:7:W vs. u[p][i + 1]@38:50:R
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
	// Race: overlapping slices - should be b[n:2*n] but is b[n-2:2*n-2]
	u := [2][]float64{b[0:n], b[n-2 : 2*n-2]} // Overlapping ranges!

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
				// Race: u[0] and u[1] overlap, causing concurrent access to same memory
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
