/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Example with loop-carried data dependence at the outer level loop.
//The inner level loop can be parallelized. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, j, n, m;
    float** b;

    n = 100;
    m = 100;

    // Allocate 2D array
    b = (float**)malloc(n * sizeof(float*));
    for (i = 0; i < n; i++) {
        b[i] = (float*)malloc(m * sizeof(float));
    }

    for (i = 0; i < n; i++) {
        for (j = 0; j < m; j++) {
            b[i][j] = (i + 1) * (j + 1);  // Adjust for 0-based indexing
        }
    }

    for (i = 1; i < n; i++) {
        #pragma omp parallel for
        for (j = 1; j < m; j++) {
            b[i][j] = b[i-1][j-1];
        }
    }

    // Free 2D array
    for (i = 0; i < n; i++) {
        free(b[i]);
    }
    free(b);

    return 0;
} 