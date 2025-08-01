/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A kernel with two level parallelizable loop with reduction:
//if reduction(+:sum) is missing, there is race condition.
//Data race pairs: we allow multiple pairs to preserve the pattern.
//  getSum@37:13:W vs. getSum@37:13:W
//  getSum@37:13:W vs. getSum@37:22:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, j, len;
    float temp, getSum;
    float** u;

    len = 100;
    getSum = 0.0;

    // Allocate 2D array
    u = (float**)malloc(len * sizeof(float*));
    for (i = 0; i < len; i++) {
        u[i] = (float*)malloc(len * sizeof(float));
    }

    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            u[i][j] = 0.5;
        }
    }

    #pragma omp parallel for private(temp, i, j)
    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            temp = u[i][j];
            getSum = getSum + temp * temp;
        }
    }

    printf("sum = %f\n", getSum);

    // Free 2D array
    for (i = 0; i < len; i++) {
        free(u[i]);
    }
    free(u);

    return 0;
} 