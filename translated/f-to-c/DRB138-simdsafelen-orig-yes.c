/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The safelen(2) clause safelen(2)@22:16 guarantees that the vector code is safe for vectors up to 2 (inclusive).
//In the loop, m can be 2 or more for the correct execution. If the value of m is less than 2,
//the behavior is undefined. Data Race Pair: b[i]@24:9:W vs. b[i-m]@24:16:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, m, n;
    float b[4];

    m = 1;  // m < 2 violates safelen(2) requirements
    n = 4;

    #pragma omp simd safelen(2)
    for (i = m + 1; i <= n; i++) {  // Adjust for 1-based to 0-based indexing
        b[i - 1] = b[i - 1 - m] - 1.0f;  // Data race: m=1 < safelen=2
    }

    printf("%f\n", b[2]);  // Adjust for 0-based indexing

    return 0;
} 