/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB100-task-reference-orig-no.cpp

Description: Cover the implicitly determined rule: In an orphaned task generating construct,
formal arguments passed by reference are firstprivate.
This requires OpenMP 4.5 to work.

In Go, we simulate this using closures that capture values by copy.
*/

package main

import (
	"fmt"
	"sync"
)

const MYLEN = 100

var a [MYLEN]int

func genTask(i int, wg *sync.WaitGroup) {
	// In Go, we use closures to simulate task generation
	// The value of i is captured by value (simulating firstprivate behavior)
	wg.Add(1)
	go func(taskI int) { // Explicitly copy i to simulate firstprivate
		defer wg.Done()
		a[taskI] = taskI + 1
	}(i) // Pass i by value to ensure task gets private copy
}

func main() {
	var wg sync.WaitGroup

	// Simulate parallel + single construct
	var taskGeneratorWg sync.WaitGroup
	taskGeneratorWg.Add(1)

	go func() {
		defer taskGeneratorWg.Done()
		// Single thread generates all tasks
		for i := 0; i < MYLEN; i++ {
			genTask(i, &wg)
		}
	}()

	// Wait for task generator to finish
	taskGeneratorWg.Wait()

	// Wait for all tasks to complete
	wg.Wait()

	// Correctness checking
	for i := 0; i < MYLEN; i++ {
		if a[i] != i+1 {
			fmt.Printf("warning: a[%d] = %d, not expected %d\n", i, a[i], i+1)
		}
	}
}
