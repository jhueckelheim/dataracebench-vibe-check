/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Only the outmost loop can be parallelized in this program.
//The inner loop has true dependence.
//Data race pair: b[i][j]@29:13:W vs. b[i][j-1]@29:22:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, j, n, m, len;
    float** b;

    len = 1000;
    n = len;
    m = len;

    // Allocate 2D array
    b = (float**)malloc(len * sizeof(float*));
    for (i = 0; i < len; i++) {
        b[i] = (float*)malloc(len * sizeof(float));
    }

    for (i = 0; i < n; i++) {
        #pragma omp parallel for
        for (j = 1; j < m; j++) {
            b[i][j] = b[i][j-1];
        }
    }

    printf("b(500,500) = %f\n", b[499][499]);

    // Free 2D array
    for (i = 0; i < len; i++) {
        free(b[i]);
    }
    free(b);

    return 0;
} 