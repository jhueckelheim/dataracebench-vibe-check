/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//micro-bench equivalent to the restrict keyword in C-99 in F95. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

void foo(int n, int** a, int** b, int** c, int** d)
{
    int i;

    *a = (int*)malloc(n * sizeof(int));
    *b = (int*)malloc(n * sizeof(int));
    *c = (int*)malloc(n * sizeof(int));
    *d = (int*)malloc(n * sizeof(int));

    for (i = 0; i < n; i++) {
        (*b)[i] = i + 1;  // Adjust for 0-based indexing
        (*c)[i] = i + 1;  // Adjust for 0-based indexing
    }

    #pragma omp parallel for
    for (i = 0; i < n; i++) {
        (*a)[i] = (*b)[i] + (*c)[i];
    }

    if ((*a)[499] != 1000) {  // Adjust for 0-based indexing
        printf("%d\n", (*a)[499]);
    }

    free(*a);
    free(*b);
    free(*c);
    free(*d);
    
    *a = NULL;
    *b = NULL;
    *c = NULL;
    *d = NULL;
}

int main()
{
    int n = 1000;
    int* a = NULL;
    int* b = NULL;
    int* c = NULL;
    int* d = NULL;

    foo(n, &a, &b, &c, &d);

    return 0;
} 