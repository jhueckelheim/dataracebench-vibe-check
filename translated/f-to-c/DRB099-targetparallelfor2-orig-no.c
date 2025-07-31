/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//use of omp target + map + array sections derived from pointers. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

float foo(double* a, double* b, int N)
{
    int i;

    #pragma omp target map(to:a[0:N]) map(from:b[0:N])
    #pragma omp parallel for
    for (i = 0; i < N; i++) {
        b[i] = a[i] * (double)(i + 1);  // Adjust for 0-based indexing
    }

    return 0.0;  // Fortran function returns but value not used
}

int main()
{
    int i, len;
    double* a;
    double* b;
    float x;

    len = 1000;

    a = (double*)malloc(len * sizeof(double));
    b = (double*)malloc(len * sizeof(double));

    for (i = 0; i < len; i++) {
        a[i] = (double)(i + 1) / 2.0;  // Adjust for 0-based indexing
        b[i] = 0.0;
    }

    x = foo(a, b, len);
    printf("b(50) = %f\n", b[49]);  // Adjust for 0-based indexing

    free(a);
    free(b);

    return 0;
} 