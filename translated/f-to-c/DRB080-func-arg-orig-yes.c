/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

// argument pass-by-reference
// its data-sharing attribute is the same as its actual argument's. i and q are shared.
// Data race pair: q@15:5:W vs. q@15:5:W

#include <omp.h>
#include <stdio.h>

void f1(int* q)
{
    *q = *q + 1;
}

int main()
{
    int i;

    i = 0;

    #pragma omp parallel
    f1(&i);

    printf("i = %d\n", i);

    return 0;
} 