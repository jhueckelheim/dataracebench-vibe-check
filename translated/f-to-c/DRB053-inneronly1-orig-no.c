/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Example with loop-carried data dependence at the outer level loop.
//But the inner level loop can be parallelized.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, j;
    float** a;

    // Allocate 2D array
    a = (float**)malloc(20 * sizeof(float*));
    for (i = 0; i < 20; i++) {
        a[i] = (float*)malloc(20 * sizeof(float));
    }

    for (i = 0; i < 20; i++) {
        for (j = 0; j < 20; j++) {
            a[i][j] = 0.0;
        }
    }

    for (i = 0; i < 19; i++) {
        #pragma omp parallel for
        for (j = 0; j < 20; j++) {
            a[i][j] = a[i][j] + a[i+1][j];
        }
    }

    // Free 2D array
    for (i = 0; i < 20; i++) {
        free(a[i]);
    }
    free(a);

    return 0;
} 