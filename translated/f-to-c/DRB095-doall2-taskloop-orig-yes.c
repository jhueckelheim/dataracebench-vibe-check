/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two-dimensional array computation:
//Only one loop is associated with omp taskloop.
//The inner loop's loop iteration variable will be shared if it is shared in the enclosing context.
//Data race pairs (we allow multiple ones to preserve the pattern):
//  Write_set = {j@36:20 (implicit step +1)}
//  Read_set = {j@36:20, j@37:35}
//  Any pair from Write_set vs. Write_set  and Write_set vs. Read_set is a data race pair.

//need to run with large thread number and large num of iterations.

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
            #pragma omp taskloop
            for (i = 0; i < len; i++) {
                for (j = 0; j < len; j++) {  // j is shared, causes data race
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