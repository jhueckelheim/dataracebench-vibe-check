/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This example is from DRACC by Adrian Schmitz et al.
//Concurrent access on a counter with no lock with simd. Atomicity Violation. Intra Region.
//Data Race Pairs: var@29:13:W vs. var@29:13:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int var[8];
    int i, j;

    for (i = 0; i < 8; i++) {
        var[i] = 0;
    }

    #pragma omp target map(tofrom:var[0:8]) device(0)
    #pragma omp teams num_teams(1) thread_limit(1048)
    #pragma omp distribute parallel for
    for (i = 0; i < 20; i++) {
        #pragma omp simd
        for (j = 0; j < 8; j++) {
            var[j] = var[j] + 1;  // Data race: concurrent SIMD access without protection
        }
    }

    printf("%d\n", var[7]);  // Adjust for 0-based indexing

    return 0;
} 