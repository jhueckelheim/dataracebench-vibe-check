/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//There is a loop-carried true dependence within the outer level loop.
//Data race pair: b[i][j]@31:13:W vs. b[i-1][j-1]@31:22:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, j, n, m;
    float** b;

    n = 1000;
    m = 1000;

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