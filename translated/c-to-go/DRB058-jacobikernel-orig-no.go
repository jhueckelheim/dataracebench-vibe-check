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
Two parallel for loops within one single parallel region,
combined with private() and reduction().
*/
package main

import (
	"fmt"
	"math"
	"sync"
)

const MSIZE = 200

var n = MSIZE
var m = MSIZE
var mits = 1000
var tol = 0.0000000001
var relax = 1.0
var alpha = 0.0543
var u [MSIZE][MSIZE]float64
var f [MSIZE][MSIZE]float64
var uold [MSIZE][MSIZE]float64
var dx, dy float64

func initialize() {
	dx = 2.0 / float64(n-1)
	dy = 2.0 / float64(m-1)

	// Initialize initial condition and RHS
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			xx := int(-1.0 + dx*float64(i-1)) // -1 < x < 1
			yy := int(-1.0 + dy*float64(j-1)) // -1 < y < 1
			u[i][j] = 0.0
			f[i][j] = -1.0*alpha*(1.0-float64(xx*xx))*(1.0-float64(yy*yy)) -
				2.0*(1.0-float64(xx*xx)) - 2.0*(1.0-float64(yy*yy))
		}
	}
}

func jacobi() {
	omega := relax
	dx = 2.0 / float64(n-1)
	dy = 2.0 / float64(m-1)

	ax := 1.0 / (dx * dx)                   // X-direction coef
	ay := 1.0 / (dy * dy)                   // Y-direction coef
	b := -2.0/(dx*dx) - 2.0/(dy*dy) - alpha // Central coeff

	error := 10.0 * tol
	k := 1

	for k <= mits {
		error = 0.0

		// Copy new solution into old in parallel
		var wg1 sync.WaitGroup
		for i := 0; i < n; i++ {
			wg1.Add(1)
			go func(i int) {
				defer wg1.Done()
				for j := 0; j < m; j++ {
					uold[i][j] = u[i][j]
				}
			}(i)
		}
		wg1.Wait()

		// Compute residuals in parallel with reduction
		var mu sync.Mutex
		var wg2 sync.WaitGroup
		for i := 1; i < n-1; i++ {
			wg2.Add(1)
			go func(i int) {
				defer wg2.Done()
				localError := 0.0
				for j := 1; j < m-1; j++ {
					resid := (ax*(uold[i-1][j]+uold[i+1][j]) +
						ay*(uold[i][j-1]+uold[i][j+1]) +
						b*uold[i][j] - f[i][j]) / b

					u[i][j] = uold[i][j] - omega*resid
					localError += resid * resid
				}
				// Reduction operation
				mu.Lock()
				error += localError
				mu.Unlock()
			}(i)
		}
		wg2.Wait()

		// Error check
		k++
		error = math.Sqrt(error) / float64(n*m)
	}

	fmt.Printf("Total Number of Iterations:%d\n", k)
	fmt.Printf("Residual:%E\n", error)
}

func main() {
	initialize()
	jacobi()
}
