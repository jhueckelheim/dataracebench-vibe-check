/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This example is based on one code snippet extracted from a paper:
//Ma etc. Symbolic Analysis of Concurrency Errors in OpenMP Programs, ICPP 2013
//
//Explicit barrier to counteract nowait. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, error, len, b;
    int* a;

    len = 1000;
    b = 5;
    a = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;  // Adjust for 0-based indexing
    }

    #pragma omp parallel shared(b, error)
    {
        #pragma omp for nowait
        for (i = 0; i < len; i++) {
            a[i] = b + a[i] * 5;
        }
    }

    #pragma omp barrier
    #pragma omp single
    error = a[8] + 1;  // Adjust for 0-based indexing

    printf("error = %d\n", error);

    free(a);

    return 0;
} 