/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The safelen(2) clause safelen(2)@23:16 guarantees that the vector code is safe for vectors up to 2 (inclusive).
//In the loop, m can be 2 or more for the correct execution. If the value of m is less than 2,
//the behavior is undefined. No Data Race in b[i]@25:9 assignment.

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, m, n;
    float b[4];

    m = 2;
    n = 4;

    #pragma omp simd safelen(2)
    for (i = m + 1; i <= n; i++) {  // Adjust for 1-based to 0-based indexing
        b[i - 1] = b[i - 1 - m] - 1.0f;  // Safe due to m >= 2
    }

    printf("%f\n", b[2]);  // Adjust for 0-based indexing

    return 0;
} 