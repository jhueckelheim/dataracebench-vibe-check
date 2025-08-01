/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The outer loop has a loop-carried true dependence.
//Data race pair: b[i][j]@56:13:W vs. b[i-1][j-1]@56:22:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, j, n, m, len;
    float** b;

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

    n = len;
    m = len;

    // Allocate 2D array
    b = (float**)malloc(n * sizeof(float*));
    for (i = 0; i < n; i++) {
        b[i] = (float*)malloc(m * sizeof(float));
    }

    for (i = 0; i < n; i++) {
        for (j = 0; j < m; j++) {
            b[i][j] = 0.5;
        }
    }

    #pragma omp parallel for private(j)
    for (i = 1; i < n; i++) {
        for (j = 1; j < m; j++) {
            b[i][j] = b[i-1][j-1];
        }
    }

    printf("b(500,500) = %f\n", b[499][499]);

    // Free 2D array
    for (i = 0; i < n; i++) {
        free(b[i]);
    }
    free(b);

    return 0;
} 