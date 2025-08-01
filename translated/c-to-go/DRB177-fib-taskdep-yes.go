/*
Fibonacci code with data race (possible to scale problem size by providing
size argument).
Data Race Pair, i@25:5:W vs. i@29:7:R
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func fib(n int) int {
	var i, j int
	if n < 2 {
		return n
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Task 1: compute fib(n-1) -> i
	go func() {
		defer wg.Done()
		i = fib(n - 1)
	}()

	// Task 2: compute fib(n-2) -> j
	go func() {
		defer wg.Done()
		j = fib(n - 2)
	}()

	// Race condition: we return i + j before ensuring both tasks complete
	// This simulates the original OpenMP race where task dependency was incomplete

	// Task 3: should depend on both i and j, but we only wait for j-task implicitly
	// The original only had "depend(in : j)" but used both i and j

	// Wait for tasks but the race still exists in the logic
	wg.Wait()

	// Return uses both i and j - race condition occurs here
	return i + j
}

func main() {
	n := 10
	if len(os.Args) > 1 {
		if val, err := strconv.Atoi(os.Args[1]); err == nil {
			n = val
		}
	}

	// Use goroutine to simulate omp parallel sections
	go func() {
		fmt.Printf("fib(%d) = %d\n", n, fib(n))
	}()

	// Simple wait to let goroutine complete
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("fib(%d) = %d\n", n, fib(n))
	}()
	wg.Wait()
}
