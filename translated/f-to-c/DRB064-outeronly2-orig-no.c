/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Only the outmost loop can be parallelized.
//The inner loop has loop carried true data dependence.
//However, the loop is not parallelized so no race condition.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

void foo()
{
    int i, j, n, m, len;
    float** b;

    len = 100;

    // Allocate 2D array
    b = (float**)malloc(len * sizeof(float*));
    for (i = 0; i < len; i++) {
        b[i] = (float*)malloc(len * sizeof(float));
    }

    n = len;
    m = len;

    #pragma omp parallel for private(j)
    for (i = 0; i < n; i++) {
        for (j = 1; j < m; j++) {
            b[i][j] = b[i][j-1];
        }
    }

    // Free 2D array
    for (i = 0; i < len; i++) {
        free(b[i]);
    }
    free(b);
}

int main()
{
    foo();
    return 0;
} 