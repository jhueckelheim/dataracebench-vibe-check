/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two-dimensional array computation:
//ordered(2) is used to associate two loops with omp for.
//The corresponding loop iteration variables are private.
//
//ordered(n) is an OpenMP 4.5 addition. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global variable (from module)
int** a;

int main()
{
    int len, i, j;
    
    len = 100;

    // Allocate 2D array
    a = (int**)malloc(len * sizeof(int*));
    for (i = 0; i < len; i++) {
        a[i] = (int*)malloc(len * sizeof(int));
    }

    #pragma omp parallel for ordered(2)
    for (i = 0; i < len; i++) {
        for (j = 0; j < len; j++) {
            a[i][j] = a[i][j] + 1;
            #pragma omp ordered depend(sink:i-1,j) depend(sink:i,j-1)
            printf("test i = %d  j = %d\n", i + 1, j + 1);  // Adjust for 1-based output
            #pragma omp ordered depend(source)
        }
    }

    // Free 2D array
    for (i = 0; i < len; i++) {
        free(a[i]);
    }
    free(a);

    return 0;
} 