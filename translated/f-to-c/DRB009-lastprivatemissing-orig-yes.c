/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This loop has loop-carried output-dependence due to x=... at line 21.
//The problem can be solved by using lastprivate(x).
//Data race pair: x@21:9:W vs. x@21:9:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, x, len;
    len = 10000;

    #pragma omp parallel for private(i)
    for (i = 0; i <= len; i++) {
        x = i;
    }

    printf("x = %d\n", x);

    return 0;
} 