/*
Copyright (c) 2017, Lawrence Livermore National Security, LLC.
Produced at the Lawrence Livermore National Laboratory
Written by Chunhua Liao, Pei-Hung Lin, Joshua Asplund,
Markus Schordan, and Ian Karlin
(email: liao6@llnl.gov, lin32@llnl.gov, asplund1@llnl.gov,
schordan1@llnl.gov, karlin1@llnl.gov)
LLNL-CODE-732144
All rights reserved.

This file is part of DataRaceBench. For details, see
https://github.com/LLNL/dataracebench. Please also see the LICENSE file
for our additional BSD notice.

Redistribution and use in source and binary forms, with
or without modification, are permitted provided that the following
conditions are met:

* Redistributions of source code must retain the above copyright
  notice, this list of conditions and the disclaimer below.

* Redistributions in binary form must reproduce the above copyright
  notice, this list of conditions and the disclaimer (as noted below)
  in the documentation and/or other materials provided with the
  distribution.

* Neither the name of the LLNS/LLNL nor the names of its contributors
  may be used to endorse or promote products derived from this
  software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL LAWRENCE LIVERMORE NATIONAL
SECURITY, LLC, THE U.S. DEPARTMENT OF ENERGY OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY,
OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
Two-dimensional array computation:
Demonstrating explicit variable scoping - private variables in goroutines, shared array data
*/
package main

import (
	"runtime"
	"sync"
)

var a [100][100]int
var b [100][100]int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	numThreads := runtime.NumCPU()

	// First parallel region - explicit private variables (i, j are local to each goroutine)
	wg.Add(numThreads)
	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Private variables i, j for this goroutine (equivalent to default(none) private(i,j))
			chunkSize := 100 / numThreads
			start := threadID * chunkSize
			end := start + chunkSize
			if threadID == numThreads-1 {
				end = 100
			}

			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					a[i][j] = a[i][j] + 1 // Shared array a, private loop variables
				}
			}
		}(t)
	}
	wg.Wait()

	// Second parallel region - shared by default, but private loop variables
	wg.Add(numThreads)
	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Private variables i, j for this goroutine (equivalent to default(shared) private(i,j))
			chunkSize := 100 / numThreads
			start := threadID * chunkSize
			end := start + chunkSize
			if threadID == numThreads-1 {
				end = 100
			}

			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					b[i][j] = b[i][j] + 1 // Shared array b, private loop variables
				}
			}
		}(t)
	}
	wg.Wait()
}
