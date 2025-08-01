/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB102-copyprivate-orig-no.c

Description: threadprivate+copyprivate: no data races
The copyprivate clause broadcasts values from one thread to all other threads.
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	numThreads := 4

	// Channel to broadcast copyprivate values
	type CopyPrivateData struct {
		x float32
		y int
	}
	broadcastChan := make(chan CopyPrivateData, 1)

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			// Each thread has its own private x and y (simulates threadprivate)
			var x float32 = 0.0
			var y int = 0

			// Simulate single copyprivate - only thread 0 sets values
			if threadID == 0 {
				x = 1.0
				y = 1
				// Broadcast values to all threads (copyprivate)
				broadcastChan <- CopyPrivateData{x: x, y: y}
			} else {
				// Other threads receive broadcasted values
				data := <-broadcastChan
				x = data.x
				y = data.y
				// Put it back for next thread (if any)
				broadcastChan <- data
			}

			// Each thread prints its private values (should all be 1.0, 1)
			fmt.Printf("Thread %d: x=%f y=%d\n", threadID, x, y)
		}(t)
	}

	wg.Wait()
}
