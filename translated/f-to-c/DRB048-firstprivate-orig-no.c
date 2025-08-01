/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Example use of firstprivate(). No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global array (from module)
int* a;

void foo(int* a, int n, int g)
{
    int i;

    #pragma omp parallel for firstprivate(g)
    for (i = 0; i < n; i++) {
        a[i] = a[i] + g;
    }
}

int main()
{
    a = (int*)malloc(100 * sizeof(int));
    
    foo(a, 100, 7);
    
    printf("%d\n", a[49]);
    
    free(a);
    
    return 0;
} 