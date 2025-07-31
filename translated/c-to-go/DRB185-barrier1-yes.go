/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

wrong 2-thread flag barrier using busy-waits, race
Data race pair: x@85:9:W vs. x@92:9:W
*/

package main

import (
	"fmt"
	"sync"
)

type flagT bool

var (
	f0, f1 flagT
	n      = 2
	x      = 1
	mutex  sync.Mutex
)

func initialize(f *flagT) {
	*f = false
}

func raise(f *flagT) {
	mutex.Lock()
	defer mutex.Unlock()
	if *f != false {
		panic("Assertion failed: flag should be false")
	}
	*f = true
}

func lower(f *flagT) {
	done := false
	for !done {
		mutex.Lock()
		if *f == true {
			*f = false
			done = true
		}
		mutex.Unlock()
	}
}

func myBarrier(tid int) {
	// This is the faulty barrier - each thread only waits on its own flag
	// Thread 0 raises f0 and waits for f0, Thread 1 raises f1 and waits for f1
	// This creates no synchronization between threads!
	if tid == 0 {
		raise(&f0)
		lower(&f0) // Race: waits for own flag, not the other thread
	} else if tid == 1 {
		raise(&f1)
		lower(&f1) // Race: waits for own flag, not the other thread
	}
}

func main() {
	initialize(&f0)
	initialize(&f1)

	var wg sync.WaitGroup
	wg.Add(2)

	// Thread 0
	go func() {
		defer wg.Done()
		tid := 0

		for i := 0; i < n; i++ {
			fmt.Printf("Thread %d: phase 1, i=%d, x=%d\n", tid, i, x)
			myBarrier(tid)
			x = 0 // Race: no proper synchronization
			myBarrier(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			myBarrier(tid)
			myBarrier(tid)
		}
	}()

	// Thread 1
	go func() {
		defer wg.Done()
		tid := 1

		for i := 0; i < n; i++ {
			fmt.Printf("Thread %d: phase 1, i=%d, x=%d\n", tid, i, x)
			myBarrier(tid)
			myBarrier(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			myBarrier(tid)
			x = 1 // Race: no proper synchronization
			myBarrier(tid)
		}
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
