/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Due to distribute parallel for simd directive at line 23, there is a data race at line 25.
//Data Rae Pairs, var@25:9:W vs. var@25:15:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int var[100];
    int i;

    for (i = 0; i < 100; i++) {
        var[i] = 1;
    }

    #pragma omp target map(tofrom:var[0:100]) device(0)
    #pragma omp teams distribute parallel for simd safelen(16)
    for (i = 16; i < 100; i++) {  // Adjust for 0-based indexing
        var[i] = var[i - 16] + 1;  // Data race: SIMD with true dependency
    }

    printf("%d\n", var[97]);  // Adjust for 0-based indexing

    return 0;
} 