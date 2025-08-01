/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This program has data races due to true dependence within the loop at 63.
//Data race pair: a[i+1]@27:9:W vs. a[i]@27:18:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, len;
    int* a;

    len = 100;
    a = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;
    }

    #pragma omp parallel for
    for (i = 0; i < len-1; i++) {
        a[i+1] = a[i] + 1;
    }

    printf("a(50) = %d\n", a[49]);

    free(a);

    return 0;
} 