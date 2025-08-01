/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//tmp should be annotated as private to avoid race condition.
//Data race pairs: tmp@28:9:W vs. tmp@29:16:R
//                 tmp@28:9:W vs. tmp@28:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, tmp, len;
    int* a;

    len = 100;
    a = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;
    }

    #pragma omp parallel for
    for (i = 0; i < len; i++) {
        tmp = a[i] + i + 1;  // Adjust for 0-based indexing
        a[i] = tmp;
    }

    printf("a(50) = %d\n", a[49]);

    free(a);

    return 0;
} 