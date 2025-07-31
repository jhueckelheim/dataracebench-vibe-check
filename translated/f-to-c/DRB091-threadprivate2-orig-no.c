/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A file-scope variable used within a function called by a parallel region.
//Use threadprivate to avoid data races.
//This is the case for a variable referenced within a construct. No data race pairs.

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
int sum0 = 0;
int sum1 = 0;

#pragma omp threadprivate(sum0)

int main()
{
    int len, i, sum;
    
    len = 1000;
    sum = 0;

    #pragma omp parallel copyin(sum0)
    {
        #pragma omp for
        for (i = 1; i <= len; i++) {
            sum0 = sum0 + i;
        }
        
        #pragma omp critical
        sum = sum + sum0;
    }

    for (i = 1; i <= len; i++) {
        sum1 = sum1 + i;
    }

    printf("sum = %d sum1 = %d\n", sum, sum1);

    return 0;
} 