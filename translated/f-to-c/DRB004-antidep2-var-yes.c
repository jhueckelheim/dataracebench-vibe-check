/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two nested loops with loop-carried anti-dependence on the outer level.
//This is a variable-length array version in F95.
//Data race pair: a[i][j]@55:13:W vs. a[i+1][j]@55:31:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, j, len;
    float** a;
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

    // Allocate 2D array
    a = (float**)malloc(len * sizeof(float*));
    for (i = 0; i < len; i++) {
        a[i] = (float*)malloc(len * sizeof(float));
    }

    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            a[i][j] = 0.5;
        }
    }

    #pragma omp parallel for private(j)
    for (i = 0; i < len-1; i++) {
        for (j = 0; j < len; j++) {
            a[i][j] = a[i][j] + a[i+1][j];
        }
    }

    printf("a(10,10) = %f\n", a[9][9]);

    // Free 2D array
    for (i = 0; i < len; i++) {
        free(a[i]);
    }
    free(a);

    return 0;
} 