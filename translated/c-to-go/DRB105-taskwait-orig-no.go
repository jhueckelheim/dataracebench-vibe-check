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

/* This is a program based on a test contributed by Yizi Gu@Rice Univ.
 * Classic Fibonacci calculation using goroutines with proper synchronization. No data races.
 * */
package main

import (
	"fmt"
	"runtime"
)

const input = 30

func fib(n uint) int {
	if n < 2 {
		return int(n)
	}

	// Channels to receive results from goroutines
	iChan := make(chan int)
	jChan := make(chan int)

	// Start two goroutines for recursive calls
	go func() {
		iChan <- fib(n - 1)
	}()

	go func() {
		jChan <- fib(n - 2)
	}()

	// Wait for both results (equivalent to taskwait)
	i := <-iChan
	j := <-jChan

	return i + j
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	result := fib(input)

	fmt.Printf("Fib(%d)=%d\n", input, result)
	if result != 832040 {
		panic("Wrong result!")
	}
}
