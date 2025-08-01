/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A file-scope variable used within a function called by a parallel region.
//No threadprivate is used to avoid data races.
//
//Data race pairs  sum@39:13:W vs. sum@39:19:R
//                 sum@39:13:W vs. sum@39:13:W

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
long long sum0 = 0;
long long sum1 = 0;

void foo(long long i)
{
    sum0 = sum0 + i;  // Data race on global sum0
}

int main()
{
    long long i, sum;
    sum = 0;

    #pragma omp parallel
    {
        #pragma omp for
        for (i = 1; i <= 1001; i++) {
            foo(i);
        }
        
        #pragma omp critical
        sum = sum + sum0;
    }

    for (i = 1; i <= 1001; i++) {
        sum1 = sum1 + i;
    }

    printf("sum = %lld sum1 = %lld\n", sum, sum1);

    return 0;
} 