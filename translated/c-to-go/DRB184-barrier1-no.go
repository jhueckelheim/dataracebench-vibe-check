/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

2-thread flag barrier using busy-wait loops and critical, no race.
*/

package main

import (
	"fmt"
	"sync"
)

type flagT bool

var (
	f0 flagT
	f1 flagT
	n  = 100
	x  = 1
)

var mutex sync.Mutex

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
	if tid == 0 {
		raise(&f0)
		lower(&f1)
	} else if tid == 1 {
		lower(&f0)
		raise(&f1)
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

		// Initial barrier equivalent
		// (in Go we don't need explicit barrier here)

		for i := 0; i < n; i++ {
			fmt.Printf("Thread %d: phase 1, i=%d, x=%d\n", tid, i, x)
			if x != 1 {
				panic("Assertion failed: x should be 1")
			}
			myBarrier(tid)
			x = 0
			myBarrier(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			if x != 0 {
				panic("Assertion failed: x should be 0")
			}
			myBarrier(tid)
			myBarrier(tid)
		}
	}()

	// Thread 1
	go func() {
		defer wg.Done()
		tid := 1

		// Initial barrier equivalent
		// (in Go we don't need explicit barrier here)

		for i := 0; i < n; i++ {
			fmt.Printf("Thread %d: phase 1, i=%d, x=%d\n", tid, i, x)
			if x != 1 {
				panic("Assertion failed: x should be 1")
			}
			myBarrier(tid)
			myBarrier(tid)
			fmt.Printf("Thread %d: phase 3, i=%d, x=%d\n", tid, i, x)
			if x != 0 {
				panic("Assertion failed: x should be 0")
			}
			myBarrier(tid)
			x = 1
			myBarrier(tid)
		}
	}()

	wg.Wait()
	fmt.Printf("Done: x=%d\n", x)
}
