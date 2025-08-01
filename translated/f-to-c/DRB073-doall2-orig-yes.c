/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two-dimensional array computation using loops: missing private(j).
//References to j in the loop cause data races.
//Data race pairs (we allow multiple ones to preserve the pattern):
//  Write_set = {j@28:12} (implicit step by +1)
//  Read_set = {j@29:17, j@29:26, j@28:12} (implicit step by +1)
//  Any pair from Write_set vs. Write_set  and Write_set vs. Read_set is a data race pair.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, j, len;
    int** a;

    len = 100;

    // Allocate 2D array
    a = (int**)malloc(len * sizeof(int*));
    for (i = 0; i < len; i++) {
        a[i] = (int*)malloc(len * sizeof(int));
    }

    #pragma omp parallel for
    for (i = 0; i < 100; i++) {
        for (j = 0; j < 100; j++) {  // Missing private(j) - causes data race
            a[i][j] = a[i][j] + 1;
        }
    }

    // Free 2D array
    for (i = 0; i < len; i++) {
        free(a[i]);
    }
    free(a);

    return 0;
} 