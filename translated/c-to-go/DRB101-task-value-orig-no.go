/*
DataRaceBench translation to Go
Copyright (c) 2017, Lawrence Livermore National Security, LLC.

This is a translation of DRB101-task-value-orig-no.cpp

Description: Cover an implicitly determined rule: In a task generating construct,
a variable without applicable rules is firstprivate.

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
	// The value of i is naturally captured by value (simulating firstprivate behavior)
	wg.Add(1)
	go func() {
		defer wg.Done()
		// i is captured by closure - acts as firstprivate
		a[i] = i + 1
	}()
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
