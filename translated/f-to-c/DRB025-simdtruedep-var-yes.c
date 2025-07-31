/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This one has race condition due to true dependence.
//But data races happen at instruction level, not thread level.
//Data race pair: a[i+1]@55:18:R vs. a[i]@55:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, len;
    int* a;
    int* b;

    len = 100;

    if (argc == 1) {
        printf("No command line arguments provided.\n");
    }

    if (argc >= 2) {
        len = atoi(argv[1]);
        if (len <= 0) {
            printf("Error, invalid integer value.\n");
        }
    }

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