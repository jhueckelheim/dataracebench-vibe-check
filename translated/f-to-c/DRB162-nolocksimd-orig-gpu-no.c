/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Reduction clause at line 23:34 will ensure there is no data race in var@27:13. No Data Race.

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
    #pragma omp distribute parallel for reduction(+:var)
    for (i = 0; i < 20; i++) {
        #pragma omp simd
        for (j = 0; j < 8; j++) {
            var[j] = var[j] + 1;  // No race: reduction clause protects access
        }
    }

    for (i = 0; i < 8; i++) {
        if (var[i] != 20) {
            printf("%d\n", var[i]);
        }
    }

    return 0;
} 