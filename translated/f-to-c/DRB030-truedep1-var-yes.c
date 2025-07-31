/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This program has data races due to true dependence within a loop.
//Data race pair: a[i+1]@51:9:W vs. a[i]@51:18:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, len;
    int* a;

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