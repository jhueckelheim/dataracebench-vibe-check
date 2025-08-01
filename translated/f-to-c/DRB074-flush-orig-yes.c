/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This benchmark is extracted from flush_nolist.1c of OpenMP Application
//Programming Interface Examples Version 4.5.0 .
//We added one critical section to make it a test with only one pair of data races.
//The data race will not generate wrong result though. So the assertion always passes.
//Data race pair:  i@37:13:W vs. i@38:15:R

#include <omp.h>
#include <stdio.h>

void f1(int* q)
{
    #pragma omp critical
    *q = 1;
    #pragma omp flush
}

int main()
{
    int i, sum;
    i = 0;
    sum = 0;

    #pragma omp parallel reduction(+:sum) num_threads(10)
    {
        f1(&i);
        sum = sum + i;
    }

    if (sum != 10) {
        printf("sum = %d\n", sum);
    }

    return 0;
} 