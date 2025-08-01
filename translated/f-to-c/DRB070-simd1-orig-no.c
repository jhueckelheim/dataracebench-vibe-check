/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//One dimension array computation with a vectorization directive. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int len, i;
    int* a;
    int* b;
    int* c;

    len = 100;

    a = (int*)malloc(len * sizeof(int));
    b = (int*)malloc(len * sizeof(int));
    c = (int*)malloc(len * sizeof(int));

    #pragma omp simd
    for (i = 0; i < len; i++) {
        a[i] = b[i] + c[i];
    }

    free(a);
    free(b);
    free(c);

    return 0;
} 