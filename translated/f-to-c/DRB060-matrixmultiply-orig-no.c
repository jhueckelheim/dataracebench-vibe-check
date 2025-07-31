/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Classic i-k-j matrix multiplication. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int N, M, K, len, i, j, l;
    float** a;
    float** b;
    float** c;

    len = 100;
    N = len;
    M = len;
    K = len;

    // Allocate 2D arrays
    a = (float**)malloc(N * sizeof(float*));
    b = (float**)malloc(M * sizeof(float*));
    c = (float**)malloc(K * sizeof(float*));
    for (i = 0; i < N; i++) {
        a[i] = (float*)malloc(M * sizeof(float));
    }
    for (i = 0; i < M; i++) {
        b[i] = (float*)malloc(K * sizeof(float));
    }
    for (i = 0; i < K; i++) {
        c[i] = (float*)malloc(N * sizeof(float));
    }

    #pragma omp parallel for private(j, l)
    for (i = 0; i < N; i++) {
        for (l = 0; l < K; l++) {
            for (j = 0; j < M; j++) {
                c[i][j] = c[i][j] + a[i][l] * b[l][j];
            }
        }
    }

    // Free 2D arrays
    for (i = 0; i < N; i++) {
        free(a[i]);
    }
    for (i = 0; i < M; i++) {
        free(b[i]);
    }
    for (i = 0; i < K; i++) {
        free(c[i]);
    }
    free(a);
    free(b);
    free(c);

    return 0;
} 