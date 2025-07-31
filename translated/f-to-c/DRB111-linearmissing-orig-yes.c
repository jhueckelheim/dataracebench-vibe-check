/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

/* loop missing the linear clause
 * Data race pair:  j@37:11:R vs. j@38:9:W
 *                  j@37:18:R vs. j@38:9:W
 */

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

    #pragma omp parallel for  // Missing linear(j:1) clause
    for (i = 0; i < len; i++) {
        c[j] = c[j] + a[i] * b[i];  // Data race on j
        j = j + 1;                  // Data race on j
    }

    printf("c(50) = %f\n", c[49]);  // Adjust for 0-based indexing

    free(a);
    free(b);
    free(c);

    return 0;
} 