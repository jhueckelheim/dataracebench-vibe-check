/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A kernel for two level parallelizable loop with reduction:
//if reduction(+:sum) is missing, there is race condition.
//Data race pairs:
//  getSum@60:13:W vs. getSum@60:13:W
//  getSum@60:13:W vs. getSum@60:22:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, j, len;
    float temp, getSum;
    float** u;

    len = 100;
    getSum = 0.0;

    if (argc == 1) {
        printf("No command line arguments provided.\n");
    }

    if (argc >= 2) {
        len = atoi(argv[1]);
        if (len <= 0) {
            printf("Error, invalid integer value.\n");
        }
    }

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