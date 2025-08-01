/*
Input dependence race: example from OMPRacer: A Scalable and Precise Static Race
Detector for OpenMP Programs
Data Race Pair, A[0]@45:7:W vs. A[i]@42:5:W
*/

package main

import (
	"os"
	"strconv"
	"sync"
)

func loadFromInput(data []int, size int) {
	for i := 0; i < size; i++ {
		data[i] = size - i
	}
}

func main() {
	N := 100

	if len(os.Args) > 1 {
		if val, err := strconv.Atoi(os.Args[1]); err == nil {
			N = val
		}
	}

	A := make([]int, N)

	loadFromInput(A, N)

	// Parallel for loop with race condition
	var wg sync.WaitGroup

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Normal assignment - each thread writes to different index
			A[idx] = idx

			// Race condition: when N > 10000, thread also writes to A[0]
			// This creates a race between the thread handling idx=0 and
			// any other thread when N > 10000
			if N > 10000 {
				A[0] = 1 // Race: multiple threads may write to A[0]
			}
		}(i)
	}

	wg.Wait()
}
