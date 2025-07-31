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
  - loop missing the linear clause
  - Data race pairs (race on j allows wrong indexing of c):
    j (read vs. write)
    j (write vs. write)
    c[j] (write vs. write due to wrong indexing)
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	const length = 100
	var a, b, c [length]float64
	var j int // Shared index variable - causes data races
	var wg sync.WaitGroup

	// Initialize arrays
	for i := 0; i < length; i++ {
		a[i] = float64(i) / 2.0
		b[i] = float64(i) / 3.0
		c[i] = float64(i) / 7.0
	}

	numThreads := runtime.NumCPU()
	wg.Add(numThreads)

	// Parallel loop with shared index variable
	for t := 0; t < numThreads; t++ {
		go func(threadID int) {
			defer wg.Done()

			// Each goroutine processes a chunk of iterations
			chunkSize := length / numThreads
			start := threadID * chunkSize
			end := start + chunkSize
			if threadID == numThreads-1 {
				end = length // Handle remainder for last thread
			}

			for i := start; i < end; i++ {
				// Data race: concurrent read of j for indexing
				if j < length {
					c[j] += a[i] * b[i] // Data race: concurrent writes to c[j] due to shared j
				}
				j++ // Data race: concurrent writes to j
			}
		}(t)
	}

	wg.Wait()

	fmt.Printf("c[50]=%f\n", c[50])
}
