/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Loop carried true dep between tmp =..  and ..= tmp.
//Data race pair: tmp@24:9:W vs. tmp@25:15:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, tmp, len;
    int* a;

    len = 100;
    tmp = 10;
    a = (int*)malloc(len * sizeof(int));

    #pragma omp parallel for
    for (i = 0; i < len; i++) {
        a[i] = tmp;
        tmp = a[i] + i + 1;  // Adjust for 0-based indexing
    }

    printf("a(50) = %d\n", a[49]);

    free(a);

    return 0;
} 