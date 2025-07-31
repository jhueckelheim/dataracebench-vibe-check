/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This one has data races due to true dependence.
//But data races happen at instruction level, not thread level.
//Data race pair: a[i+1]@32:9:W vs. a[i]@32:18:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, len;
    int* a;
    int* b;

    len = 100;

    a = (int*)malloc(len * sizeof(int));
    b = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;
        b[i] = i + 2;
    }

    #pragma omp simd
    for (i = 0; i < len-1; i++) {
        a[i+1] = a[i] + b[i];
    }

    for (i = 0; i < len; i++) {
        printf("Values for i and a(i) are: %d %d\n", i, a[i]);
    }

    free(a);
    free(b);

    return 0;
} 