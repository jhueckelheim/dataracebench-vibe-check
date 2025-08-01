/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Data race pair: a[i]@49:9:W vs. a[0]@49:16:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int len, i;
    int* a;

    len = 1000;

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

    a[0] = 2;

    #pragma omp parallel for
    for (i = 0; i < len; i++) {
        a[i] = a[i] + a[0];
    }

    printf("a(0) = %d\n", a[0]);

    free(a);

    return 0;
} 