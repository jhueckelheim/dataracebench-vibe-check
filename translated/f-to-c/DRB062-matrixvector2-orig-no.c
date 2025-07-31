/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Matrix-vector multiplication: inner level parallelization. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

void foo()
{
    int i, j, N;
    float sum;
    float** a;
    float* v;
    float* v_out;

    N = 1000;

    // Allocate arrays
    a = (float**)malloc(N * sizeof(float*));
    for (i = 0; i < N; i++) {
        a[i] = (float*)malloc(N * sizeof(float));
    }
    v = (float*)malloc(N * sizeof(float));
    v_out = (float*)malloc(N * sizeof(float));

    for (i = 0; i < N; i++) {
        sum = 0.0;
        #pragma omp parallel for reduction(+:sum)
        for (j = 0; j < N; j++) {
            sum = sum + a[i][j] * v[j];
            printf("%f\n", sum);
        }
        v_out[i] = sum;
    }

    // Free arrays
    for (i = 0; i < N; i++) {
        free(a[i]);
    }
    free(a);
    free(v);
    free(v_out);
}

int main()
{
    foo();
    return 0;
} 