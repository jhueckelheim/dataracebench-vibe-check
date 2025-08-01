/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two-dimensional array computation:
//default(none) to enforce explictly list all variables in data-sharing attribute clauses
//default(shared) to cover another option. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int len, i, j;
    double** a;
    double** b;

    len = 100;

    // Allocate 2D arrays
    a = (double**)malloc(len * sizeof(double*));
    b = (double**)malloc(len * sizeof(double*));
    for (i = 0; i < len; i++) {
        a[i] = (double*)malloc(len * sizeof(double));
        b[i] = (double*)malloc(len * sizeof(double));
    }

    #pragma omp parallel for default(none) shared(a) private(i, j)
    for (i = 0; i < 100; i++) {
        for (j = 0; j < 100; j++) {
            a[i][j] = a[i][j] + 1;
        }
    }

    #pragma omp parallel for default(shared) private(i, j)
    for (i = 0; i < 100; i++) {
        for (j = 0; j < 100; j++) {
            b[i][j] = b[i][j] + 1;
        }
    }

    printf("%f %f\n", a[49][49], b[49][49]);  // Adjust for 0-based indexing

    // Free 2D arrays
    for (i = 0; i < len; i++) {
        free(a[i]);
        free(b[i]);
    }
    free(a);
    free(b);

    return 0;
} 