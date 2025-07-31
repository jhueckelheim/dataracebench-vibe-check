/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The loop in this example cannot be parallelized.
//
//This pattern has two pair of dependencies:
//1. loop carried output dependence
// x = .. :
//
//2. loop carried true dependence due to:
//.. = x;
// x = ..;
//Data race pairs: we allow two pairs to preserve the original code pattern.
// 1. x@48:16:R vs. x@49:9:W
// 2. x@49:9:W vs. x@49:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global array (from module)
int* a;

void useGlobalArray(int len)
{
    len = 100;
    a = (int*)malloc(100 * sizeof(int));
}

int main()
{
    int len, i, x;

    len = 100;
    x = 10;

    useGlobalArray(len);

    #pragma omp parallel for
    for (i = 0; i < len; i++) {
        a[i] = x;
        x = i + 1;  // Adjust for 0-based indexing in C
    }

    printf("x = %d\n", x);

    free(a);

    return 0;
} 