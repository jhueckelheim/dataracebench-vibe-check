/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

one synchronization commented out.
enters. So race on x can occur.
Data race pair: x@104:9:W vs. x@111:9:W
*/

package main

import (
	"fmt"
	"sync"
)

var (
	l0, l1, l2 sync.Mutex
	n          = 2
	x          = 1
)

func barrierInit() {
	// Locks are initialized unlocked in Go
}

func barrierStart(tid int) {
	if tid == 0 {
		l0.Lock()
		l2.Lock()
	} else if tid == 1 {
		l1.Lock()
	}
}

func barrierStop(tid int) {
	if tid == 0 {
		l0.Unlock()
		l2.Unlock()
	} else if tid == 1 {
		l1.Unlock()
	}
}

func barrierWait(tid int) {
	// Race condition: some synchronization operations are commented out
	// This breaks the barrier's correctness
	if tid == 0 {
		l0.Unlock()
		l1.Lock()
		// l2.Unlock()  // COMMENTED OUT - breaks synchronization!
		l0.Lock()
		l1.Unlock()
		// l2.Lock()    // COMMENTED OUT - breaks synchronization!
	} else if tid == 1 {
		l0.Lock()
		l1.Unlock()
		// l2.Lock()    // COMMENTED OUT - breaks synchronization!
		l0.Unlock()
		l1.Lock()
		// l2.Unlock()  // COMMENTED OUT - breaks synchronization!
	}
}

func main() {
	barrierInit()

	var wg sync.WaitGroup
	wg.Add(2)

	// Thread 0
	go func() {
		defer wg.Done()
		tid := 0
		barrierStart(tid)

		for i := 0; i < n; i++ {
			fmt.Printf("Thread %d: phase 1, i=%d, x=%d\n", tid, i, x)
			barrierWait(tid)
			x = 0 // Race: broken barrier allows concurrent access
			barrierWait(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			barrierWait(tid)
			barrierWait(tid)
		}

		barrierStop(tid)
	}()

	// Thread 1
	go func() {
		defer wg.Done()
		tid := 1
		barrierStart(tid)

		for i := 0; i < n; i++ {
			fmt.Printf("Thread %d: phase 1, i=%d, x=%d\n", tid, i, x)
			barrierWait(tid)
			barrierWait(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			barrierWait(tid)
			x = 1 // Race: broken barrier allows concurrent access
			barrierWait(tid)
		}

		barrierStop(tid)
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
