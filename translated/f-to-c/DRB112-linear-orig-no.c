/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//omp for loop is allowed to use the linear clause, an OpenMP 4.5 addition. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int len, i, j;
    double* a;
    double* b;
    double* c;

    len = 100;
    i = 0;
    j = 0;

    a = (double*)malloc(len * sizeof(double));
    b = (double*)malloc(len * sizeof(double));
    c = (double*)malloc(len * sizeof(double));

    for (i = 0; i < len; i++) {
        a[i] = (double)(i + 1) / 2.0;  // Adjust for 0-based indexing
        b[i] = (double)(i + 1) / 3.0;
        c[i] = (double)(i + 1) / 7.0;
    }

    #pragma omp parallel for linear(j:1)  // Linear clause fixes the race
    for (i = 0; i < len; i++) {
        c[j] = c[j] + a[i] * b[i];
        j = j + 1;
    }

    // printf("c(50) = %f\n", c[49]);  // Commented out as in original

    free(a);
    free(b);
    free(c);

    return 0;
} 