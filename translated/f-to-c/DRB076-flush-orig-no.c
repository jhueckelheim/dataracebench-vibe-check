/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This benchmark is extracted from flush_nolist.1c of OpenMP
//Application Programming Interface Examples Version 4.5.0 .
//
//We privatize variable i to fix data races in the original example.
//Once i is privatized, flush is no longer needed. No data race pairs.

#include <omp.h>
#include <stdio.h>

void f1(int* q)
{
    *q = 1;
}

int main()
{
    int i, sum;
    i = 0;
    sum = 0;

    #pragma omp parallel reduction(+:sum) num_threads(10) private(i)
    {
        f1(&i);
        sum = sum + i;
    }

    if (sum != 10) {
        printf("sum = %d\n", sum);
    }

    return 0;
} 