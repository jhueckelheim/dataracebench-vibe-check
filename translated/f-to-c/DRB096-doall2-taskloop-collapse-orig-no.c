/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two-dimensional array computation:
//Two loops are associated with omp taskloop due to collapse(2).
//Both loop index variables are private.
//taskloop requires OpenMP 4.5 compilers. No data race pairs.

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

    #pragma omp parallel
    {
        #pragma omp single
        {
            #pragma omp taskloop collapse(2)
            for (i = 0; i < len; i++) {
                for (j = 0; j < len; j++) {  // Both i and j are private due to collapse(2)
                    a[i][j] = a[i][j] + 1;
                }
            }
        }
    }

    printf("a(50,50) = %d\n", a[49][49]);  // Adjust for 0-based indexing

    // Free 2D array
    for (i = 0; i < len; i++) {
        free(a[i]);
    }
    free(a);

    return 0;
} 