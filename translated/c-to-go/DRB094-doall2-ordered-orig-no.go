/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB094-doall2-ordered-orig-no.c

Description: Two-dimensional array computation:
ordered(2) is used to associate two loops with omp for.
The corresponding loop iteration variables are private.

Note: ordered(n) is an OpenMP 4.5 addition requiring dependency tracking.
In Go, we simulate this with channels to maintain order.
*/

package main

import (
	"fmt"
	"sync"
)

var a [100][100]int

func main() {
	var wg sync.WaitGroup

	// Create channels for dependency tracking (simplified ordered simulation)
	done := make([][]chan bool, 100)
	for i := 0; i < 100; i++ {
		done[i] = make([]chan bool, 100)
		for j := 0; j < 100; j++ {
			done[i][j] = make(chan bool, 1)
		}
	}

	// Process elements in ordered fashion using dependency channels
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()

				// Wait for dependencies (sink: i-1,j and i,j-1)
				if i > 0 {
					<-done[i-1][j] // Wait for (i-1,j) to complete
				}
				if j > 0 {
					<-done[i][j-1] // Wait for (i,j-1) to complete
				}

				// Do the computation
				a[i][j] = a[i][j] + 1

				// Ordered section equivalent
				fmt.Printf("test i=%d j=%d\n", i, j)

				// Signal completion (source)
				done[i][j] <- true
			}(i, j)
		}
	}

	wg.Wait()
}
