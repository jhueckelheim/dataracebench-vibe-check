/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//When if() evalutes to true, this program has data races due to true dependence within the loop at 31.
//Data race pair: a[i+1]@32:9:W vs. a[i]@32:18:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, len, rem, j;
    float u;
    double* a;

    len = 100;
    a = (double*)malloc(len * sizeof(double));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;  // Adjust for 0-based indexing
    }

    // Generate random number (simplified)
    u = (float)rand() / (float)RAND_MAX;
    j = (int)(100 * u);

    #pragma omp parallel for if ((j % 2) == 0)
    for (i = 0; i < len - 1; i++) {
        a[i + 1] = a[i] + 1;  // Data race when if condition is true
    }

    printf("a(50) = %f\n", a[49]);  // Adjust for 0-based indexing

    free(a);

    return 0;
} 