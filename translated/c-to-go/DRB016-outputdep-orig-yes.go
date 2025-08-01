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
The loop in this example cannot be parallelized.

This pattern has two pair of dependencies:
1. loop carried output dependence
 x = .. :

2. loop carried true dependence due to:
.. = x;
 x = ..;
Data race pairs: we allow two pairs to preserve the original code pattern.
 1. x@73:12:R vs. x@74:5:W
 2. x@74:5:W vs. x@74:5:W
*/

package main

import (
	"fmt"
	"sync"
)

var a = make([]int, 100)

func main() {
	length := 100
	var i int
	x := 10

	var wg sync.WaitGroup
	for i = 0; i < length; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			a[index] = x
			x = index
		}(i)
	}
	wg.Wait()
	fmt.Printf("x=%d", x)
}
