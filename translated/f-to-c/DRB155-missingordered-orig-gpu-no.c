/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//By utilizing the ordered construct @23 the execution will be sequentially consistent.
//No Data Race Pair.

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
    #pragma omp parallel for ordered
    for (i = 1; i < 100; i++) {  // Adjust for 0-based indexing
        #pragma omp ordered
        var[i] = var[i - 1] + 1;  // Sequential consistency ensured by ordered
    }

    for (i = 0; i < 100; i++) {
        if (var[i] != i + 1) {  // Adjust for 0-based indexing
            printf("Data Race Present\n");
        }
    }

    return 0;
} 