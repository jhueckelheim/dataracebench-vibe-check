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
A kernel with two level parallelizable loop with reduction:
if reduction(+:sum) is missing, there is race condition.
Data race pairs: we allow multiple pairs to preserve the pattern.
  sum@70:7:W vs. sum@70:7:W
  sum@70:7:W vs. sum@70:13:R
*/
package main

import (
	"fmt"
	"sync"
)

func main() {
	var i, j int
	var temp, sum float32
	sum = 0.0
	length := 100

	u := make([][]float32, 100)
	for i = 0; i < 100; i++ {
		u[i] = make([]float32, 100)
	}

	for i = 0; i < length; i++ {
		for j = 0; j < length; j++ {
			u[i][j] = 0.5
		}
	}

	var wg sync.WaitGroup
	for i = 0; i < length; i++ {
		wg.Add(1)
		go func(row int) {
			defer wg.Done()
			for j = 0; j < length; j++ {
				temp = u[row][j]
				sum = sum + temp*temp
			}
		}(i)
	}
	wg.Wait()

	fmt.Printf("sum = %f\n", sum)
}
