//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
// and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.
//
// SPDX-License-Identifier: (BSD-3-Clause)
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//Example of writing to a file. No data race pairs.

package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

func main() {
	var length int
	var a [1000]int
	var exist bool

	length = 1000

	for i := 1; i <= length; i++ {
		a[i-1] = i
	}

	_, err := os.Stat("mytempfile.txt")
	exist = !os.IsNotExist(err)

	var file *os.File
	var stat error
	if exist {
		file, stat = os.OpenFile("mytempfile.txt", os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, stat = os.Create("mytempfile.txt")
	}

	if stat == nil {
		//$omp parallel do
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		chunkSize := length / numCPU
		if chunkSize < 1 {
			chunkSize = 1
		}

		for start := 1; start <= length; start += chunkSize {
			end := start + chunkSize - 1
			if end > length {
				end = length
			}
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for i := start; i <= end; i++ {
					fmt.Fprintln(file, a[i-1])
				}
			}(start, end)
		}
		wg.Wait()
		//$omp end parallel do

		file.Close()
		os.Remove("mytempfile.txt")
	}
}
