/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB073-doall2-orig-yes.c

Description: Two-dimensional array computation using loops: missing private(j).
References to j in the loop cause data races.

Original Data race pairs:
  Write_set = {j@61:10, j@61:20}
  Read_set = {j@62:20, j@62:12, j@61:14, j@61:20}
  Any pair from Write_set vs. Write_set and Write_set vs. Read_set is a data race pair.
*/

package main

import (
	"fmt"
	"sync"
)

var a [100][100]int

func main() {
	var wg sync.WaitGroup
	var j int // j is shared across goroutines - this causes the data race!

	// Parallel for loop - each goroutine accesses shared variable j
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Data race: multiple goroutines read and write shared variable j
			for j = 0; j < 100; j++ {
				a[i][j] = a[i][j] + 1
			}
		}(i)
	}

	wg.Wait()

	// Print a small sample to verify execution
	fmt.Printf("Sample results: a[0][0]=%d, a[50][50]=%d, a[99][99]=%d\n",
		a[0][0], a[50][50], a[99][99])
}
