/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A file-scope variable used within a function called by a parallel region.
//No threadprivate is used to avoid data races.
//This is the case for a variable referenced within a construct.
//
//Data race pairs  sum0@34:13:W vs. sum0@34:20:R
//                 sum0@34:13:W vs. sum0@34:13:W

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
int sum0 = 0;
int sum1 = 0;

int main()
{
    int i, sum;
    
    sum = 0;
    sum0 = 0;
    sum1 = 0;

    #pragma omp parallel
    {
        #pragma omp for
        for (i = 1; i <= 1001; i++) {
            sum0 = sum0 + i;  // Data race - sum0 not threadprivate
        }
        
        #pragma omp critical
        sum = sum + sum0;
    }

    for (i = 1; i <= 1001; i++) {
        sum1 = sum1 + i;
    }

    printf("sum = %d sum1 = %d\n", sum, sum1);

    return 0;
} 