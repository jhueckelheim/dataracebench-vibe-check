/*
Iteration 0 and 1 can have conflicting writes to A[0]. But if they are scheduled to be run by
the same thread, dynamic tools may miss this.
Data Race Pair, A[0]@34:7:W vs. A[i]@31:5:W
*/

package main

import "sync"

func main() {
	N := 100
	A := make([]int, N)

	// Parallel for loop
	var wg sync.WaitGroup

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Normal assignment - each thread writes to different index
			A[idx] = idx

			// Race condition: when idx=1, thread also writes to A[0]
			// This creates a race between the thread handling idx=0 and
			// the thread handling idx=1
			if idx == 1 {
				A[0] = 1 // Race: conflicts with A[0] = 0 from idx=0 thread
			}
		}(i)
	}

	wg.Wait()
}
