/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Missing ordered directive causes data race pairs var@24:9:W vs. var@24:18:R

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
    #pragma omp teams distribute parallel for
    for (i = 1; i < 100; i++) {  // Adjust for 0-based indexing
        var[i] = var[i - 1] + 1;  // Data race: missing ordered directive
    }

    printf("%d\n", var[99]);  // Adjust for 0-based indexing

    return 0;
} 