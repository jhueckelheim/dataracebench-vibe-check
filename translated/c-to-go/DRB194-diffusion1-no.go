/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

no race, but it needs to mention that u1 and u2 are not aliased
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var (
	u1, u2 []float64
	c      = 0.2
	n      = 10
	nsteps = 10
)

func main() {
	u1 = make([]float64, n)
	u2 = make([]float64, n)

	// Initialize arrays
	for i := 1; i < n-1; i++ {
		val := rand.Float64()
		u2[i] = val
		u1[i] = val
	}
	u1[0] = 0.5
	u1[n-1] = 0.5
	u2[0] = 0.5
	u2[n-1] = 0.5

	for t := 0; t < nsteps; t++ {
		var wg sync.WaitGroup

		// Parallel for loop
		for i := 1; i < n-1; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				u2[idx] = u1[idx] + c*(u1[idx-1]+u1[idx+1]-2*u1[idx])
			}(i)
		}

		wg.Wait()

		// Proper array swapping - no aliasing
		tmp := u1
		u1 = u2
		u2 = tmp
	}

	// Print results
	for i := 0; i < n; i++ {
		fmt.Printf("%.2f ", u1[i])
	}
	fmt.Println()
}
