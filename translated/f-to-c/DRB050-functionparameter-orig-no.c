/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Arrays passed as function parameters. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global arrays (from module)
double* o1;
double* c;

void foo1(double* o1, double* c, int len)
{
    double volnew_o8;
    int i;

    #pragma omp parallel for private(volnew_o8)
    for (i = 0; i < len; i++) {
        volnew_o8 = 0.5 * c[i];
        o1[i] = volnew_o8;
    }
}

int main()
{
    o1 = (double*)malloc(100 * sizeof(double));
    c = (double*)malloc(100 * sizeof(double));

    foo1(o1, c, 100);

    free(o1);
    free(c);

    return 0;
} 