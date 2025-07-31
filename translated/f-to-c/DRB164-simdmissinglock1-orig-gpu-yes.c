/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Concurrent access of var@30:13 without acquiring locks causes atomicity violation. Data race present.
//Data Race Pairs, var@30:13:W vs. var@30:22:R

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
int var[16];
int i, j;

int main()
{
    for (i = 0; i < 16; i++) {
        var[i] = 0;
    }

    #pragma omp target map(tofrom:var[0:16]) device(0)
    #pragma omp teams distribute parallel for
    for (i = 0; i < 20; i++) {
        #pragma omp simd
        for (j = 0; j < 16; j++) {
            var[j] = var[j] + 1;  // Data race: no lock or reduction protection
        }
    }

    printf("%d\n", var[15]);  // Adjust for 0-based indexing

    return 0;
} 