/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Concurrent access of var@26:13 has no atomicity violation. No data race present.

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
    #pragma omp teams distribute parallel for reduction(+:var)
    for (i = 0; i < 20; i++) {
        #pragma omp simd
        for (j = 0; j < 16; j++) {
            var[j] = var[j] + 1;  // No race: reduction clause ensures safety
        }
    }

    for (i = 0; i < 16; i++) {
        if (var[i] != 20) {
            printf("%d %d\n", var[i], i + 1);  // Adjust for 1-based output
        }
    }

    return 0;
} 