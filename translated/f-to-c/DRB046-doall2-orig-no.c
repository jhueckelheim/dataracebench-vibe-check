/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two-dimensional array computation:
//Only one loop is associated with the omp for construct.
//The inner loop's loop iteration variable needs an explicit private() clause,
//otherwise it will be shared by default. No data race pairs.

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, j;
    int a[100][100];

    #pragma omp parallel for private(j)
    for (i = 0; i < 100; i++) {
        for (j = 0; j < 100; j++) {
            a[i][j] = a[i][j] + 1;
        }
    }

    return 0;
} 