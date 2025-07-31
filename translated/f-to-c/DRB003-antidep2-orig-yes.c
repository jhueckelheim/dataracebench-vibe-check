/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A two-level loop nest with loop carried anti-dependence on the outer level.
//Data race pair: a[i][j]@29:13:W vs. a[i+1][j]@29:31:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, j, len;
    float a[20][20];

    len = 20;

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

    return 0;
} 