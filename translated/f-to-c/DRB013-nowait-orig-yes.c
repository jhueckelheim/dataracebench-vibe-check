/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This example is extracted from a paper:
//Ma etc. Symbolic Analysis of Concurrency Errors in OpenMP Programs, ICPP 2013
//
//Some threads may finish the for loop early and execute errors = dt[10]+1
//while another thread may still be simultaneously executing
//the for worksharing region by writing to d[9], causing data races.
//
//Data race pair: a[i]@41:21:R vs. a[10]@37:17:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, error, len, b;
    int* a;

    b = 5;
    len = 1000;

    a = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;
    }

    #pragma omp parallel shared(b, error)
    {
        #pragma omp parallel
        {
            #pragma omp for nowait
            for (i = 0; i < len; i++) {
                a[i] = b + a[i] * 5;
            }
            
            #pragma omp single
            error = a[9] + 1;  // Adjust for 0-based indexing
        }
    }

    printf("error = %d\n", error);

    free(a);

    return 0;
} 