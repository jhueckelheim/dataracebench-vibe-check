/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Data race pair: a[i]@24:9:W vs. a[0]@24:16:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int len, i;
    int* a;

    len = 1000;
    a = (int*)malloc(len * sizeof(int));

    a[0] = 2;

    #pragma omp parallel for
    for (i = 0; i < len; i++) {
        a[i] = a[i] + a[0];
    }

    printf("a(500) = %d\n", a[499]);

    free(a);

    return 0;
} 