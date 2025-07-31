/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A loop with loop-carried anti-dependence.
//Data race pair: a[i+1]@25:9:W vs. a[i]@25:16:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, len;
    int a[1000];

    len = 1000;

    for (i = 0; i < len; i++) {
        a[i] = i + 1;
    }

    #pragma omp parallel for
    for (i = 0; i < len-1; i++) {
        a[i] = a[i+1] + 1;
    }

    printf("a(500)=%d\n", a[499]);
    
    return 0;
} 