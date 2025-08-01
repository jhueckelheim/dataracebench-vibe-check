/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Classic i-k-j matrix multiplication. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int len, i, j;
    int* a;
    int* b;
    int* c;

    len = 100;

    a = (int*)malloc(len * sizeof(int));
    b = (int*)malloc((len + len * len) * sizeof(int));
    c = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            b[j + i * len] = 1;
        }
        a[i] = 1;
        c[i] = 0;
    }

    #pragma omp target map(to:a[0:len],b[0:len+len*len]) map(tofrom:c[0:len]) device(0)
    #pragma omp teams distribute parallel for
    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            c[i] = c[i] + a[j] * b[j + i * len];  // No race: each thread works on different c[i]
        }
    }

    for (i = 0; i < len; i++) {
        if (c[i] != len) {
            printf("%d\n", c[i]);
        }
    }

    free(a);
    free(b);
    free(c);

    return 0;
} 