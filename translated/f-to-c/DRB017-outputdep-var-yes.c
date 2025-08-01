/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The loop in this example cannot be parallelized.
//
//Data race pairs: we allow two pairs to preserve the original code pattern.
// 1. x@50:16:R vs. x@51:9:W
// 2. x@51:9:W vs. x@51:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int len, i, x;
    int* a;

    len = 100;
    x = 10;

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

    #pragma omp parallel for
    for (i = 0; i < len; i++) {
        a[i] = x;
        x = i + 1;  // Adjust for 0-based indexing in C
    }

    printf("x = %d  a(0) = %d\n", x, a[0]);

    free(a);

    return 0;
} 