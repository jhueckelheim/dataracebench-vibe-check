/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two-dimension array computation with a vetorization directive
//collapse(2) makes simd associate with 2 loops.
//Loop iteration variables should be predetermined as lastprivate. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    double** a;
    double** b;
    double** c;
    int len, i, j;

    len = 100;

    // Allocate 2D arrays
    a = (double**)malloc(len * sizeof(double*));
    b = (double**)malloc(len * sizeof(double*));
    c = (double**)malloc(len * sizeof(double*));
    for (i = 0; i < len; i++) {
        a[i] = (double*)malloc(len * sizeof(double));
        b[i] = (double*)malloc(len * sizeof(double));
        c[i] = (double*)malloc(len * sizeof(double));
    }

    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            a[i][j] = (double)(i + 1) / 2.0;  // Adjust for 0-based indexing
            b[i][j] = (double)(i + 1) / 3.0;
            c[i][j] = (double)(i + 1) / 7.0;
        }
    }

    #pragma omp simd collapse(2)
    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            c[i][j] = a[i][j] * b[i][j];
        }
    }

    printf("c(50,50) = %f\n", c[49][49]);  // Adjust for 0-based indexing

    // Free 2D arrays
    for (i = 0; i < len; i++) {
        free(a[i]);
        free(b[i]);
        free(c[i]);
    }
    free(a);
    free(b);
    free(c);

    return 0;
} 