/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The outmost loop is be parallelized.
//But the inner level loop has out of bound access for b[i][j] when i equals to 1.
//This will case memory access of a previous column's last element.
//
//For example, an array of 4x4:
//    j=1 2 3 4
// i=1  x x x x
//   2  x x x x
//   3  x x x x
//   4  x x x x
//  outer loop: j=3,
//  inner loop: i=1
//  array element accessed b[i-1][j] becomes b[0][3], which in turn is b[4][2]
//  due to linearized column-major storage of the 2-D array.
//  This causes loop-carried data dependence between j=2 and j=3.
//
//
//Data race pair: b[i][j]@67:13:W vs. b[i-1][j]@67:22:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, j, n, m, len;
    float** b;
    len = 100;

    if (argc == 1) {
        printf("No command line arguments provided.\n");
    }

    if (argc >= 2) {
        len = atoi(argv[1]);
        if (len <= 0) {
            printf("Error, invalid integer value.\n");
        }
    }

    n = len;
    m = len;

    // Allocate 2D array
    b = (float**)malloc(n * sizeof(float*));
    for (i = 0; i < n; i++) {
        b[i] = (float*)malloc(m * sizeof(float));
    }

    #pragma omp parallel for private(i)
    for (j = 1; j < n; j++) {
        for (i = 0; i < m; i++) {
            if (i > 0) {
                b[i][j] = b[i-1][j];
            }
        }
    }

    printf("b(50,50) = %f\n", b[49][49]);

    // Free 2D array
    for (i = 0; i < n; i++) {
        free(b[i]);
    }
    free(b);

    return 0;
} 