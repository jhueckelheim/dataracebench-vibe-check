/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Race condition due to anti-dependence within a loop offloaded to accelerators.
//Data race pair: a[i]@29:13:W vs. a[i+1]@29:20:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, len;
    int* a;

    len = 1000;

    a = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;
    }

    #pragma omp target map(a[0:len])
    {
        #pragma omp parallel for
        for (i = 0; i < len-1; i++) {
            a[i] = a[i+1] + 1;
        }
    }

    for (i = 0; i < len; i++) {
        printf("Values for i and a(i) are: %d %d\n", i, a[i]);
    }

    free(a);

    return 0;
} 