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
When parallelization happens (random condition), this program has data races due to true dependence within the loop.
Data race pair: a[i+1] (write) vs. a[i] (read) - array element dependency
*/
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	const length = 100
	var a [length]int

	// Initialize array
	for i := 0; i < length; i++ {
		a[i] = i
	}

	rand.Seed(time.Now().UnixNano())

	// Conditional parallelization based on random condition
	if rand.Intn(2) == 1 {
		// Parallel execution with data race
		var wg sync.WaitGroup
		numThreads := runtime.NumCPU()
		wg.Add(numThreads)

		for t := 0; t < numThreads; t++ {
			go func(threadID int) {
				defer wg.Done()

				chunkSize := (length - 1) / numThreads
				start := threadID * chunkSize
				end := start + chunkSize
				if threadID == numThreads-1 {
					end = length - 1 // Ensure we don't go beyond array bounds
				}

				for i := start; i < end; i++ {
					a[i+1] = a[i] + 1 // Data race: a[i+1] write may conflict with a[i] read from another thread
				}
			}(t)
		}
		wg.Wait()
	} else {
		// Sequential execution - no data race
		for i := 0; i < length-1; i++ {
			a[i+1] = a[i] + 1
		}
	}

	fmt.Printf("a[50]=%d\n", a[50])
}
