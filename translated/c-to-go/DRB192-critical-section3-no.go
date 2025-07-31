/*
This is a program based on a dataset contributed by
Wenhao Wu and Stephen F. Siegel @Univ. of Delaware.

signal with busy wait loop using critical sections
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var x, s int
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	// Section 1 - Signal sender
	go func() {
		defer wg.Done()
		x = 1

		mutex.Lock()
		s = 1
		mutex.Unlock()
	}()

	// Section 2 - Signal receiver
	go func() {
		defer wg.Done()
		done := 0

		for done == 0 {
			mutex.Lock()
			if s != 0 {
				done = 1
			}
			mutex.Unlock()
		}
		x = 2
	}()

	wg.Wait()
	fmt.Printf("%d\n", x)
}
