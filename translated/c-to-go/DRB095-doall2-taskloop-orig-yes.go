/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB095-doall2-taskloop-orig-yes.c

Description: Two-dimensional array computation:
Only one loop is associated with omp taskloop.
The inner loop's loop iteration variable will be shared if it is shared in the enclosing context.

Original Data race pairs (we allow multiple ones to preserve the pattern):
  Write_set = {j@69:14, j@69:30, a[i][j]@70:11}
  Read_set = {j@69:21, j@69:30, j@70:16, a[i][j]@70:11}
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
	var j int // j is shared across all goroutines - this causes the data race!

	// Simulate OpenMP taskloop - only outer loop is parallelized
	// Inner loop variable j is shared in enclosing context
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) { // i is private (passed as parameter)
			defer wg.Done()

			// Data race: multiple goroutines read and write shared variable j
			for j = 0; j < 100; j++ {
				a[i][j] += 1
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("a[50][50]=%d\n", a[50][50])
}
