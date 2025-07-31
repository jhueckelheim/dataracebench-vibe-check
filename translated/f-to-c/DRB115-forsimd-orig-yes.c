/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This one has data races due to true dependence.
//But data races happen at both instruction and thread level.
//Data race pair: a[i+1]@31:9:W vs. a[i]@31:16:R

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
        a[i] = i + 1;  // Adjust for 0-based indexing
        b[i] = i + 2;  // Adjust for 0-based indexing
    }

    #pragma omp simd
    for (i = 0; i < len - 1; i++) {
        a[i + 1] = a[i] + b[i];  // Data race due to true dependence
    }

    printf("a(50) = %d\n", a[49]);  // Adjust for 0-based indexing

    free(a);
    free(b);

    return 0;
} 