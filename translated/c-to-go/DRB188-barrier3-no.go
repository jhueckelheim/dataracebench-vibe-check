/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

implements 2-thread reuseable barrier using 3 locks, no race.
*/

package main

import (
	"fmt"
	"sync"
)

var (
	l0, l1, l2 sync.Mutex
	n          = 100
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
	if tid == 0 {
		l0.Unlock()
		l1.Lock()
		l2.Unlock()
		l0.Lock()
		l1.Unlock()
		l2.Lock()
	} else if tid == 1 {
		l0.Lock()
		l1.Unlock()
		l2.Lock()
		l0.Unlock()
		l1.Lock()
		l2.Unlock()
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
			if x != 1 {
				panic("Assertion failed: x should be 1")
			}
			barrierWait(tid)
			x = 0
			barrierWait(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			if x != 0 {
				panic("Assertion failed: x should be 0")
			}
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
			if x != 1 {
				panic("Assertion failed: x should be 1")
			}
			barrierWait(tid)
			barrierWait(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			if x != 0 {
				panic("Assertion failed: x should be 0")
			}
			barrierWait(tid)
			x = 1
			barrierWait(tid)
		}

		barrierStop(tid)
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
