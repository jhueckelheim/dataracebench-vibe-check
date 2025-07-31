/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//tmp should be put as private to avoid race condition
//Data race pair: tmp@51:9:W vs. tmp@52:16:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, tmp, len;
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
    for (i = 0; i < len; i++) {
        tmp = a[i] + i + 1;  // Adjust for 0-based indexing
        a[i] = tmp;
    }

    printf("a(50) = %d\n", a[49]);

    free(a);

    return 0;
} 