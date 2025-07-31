/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two tasks without depend clause to protect data writes.
//i is shared for two tasks based on implicit data-sharing attribute rules.
//Data race pair: i@22:5:W vs. i@25:5:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int i;
    i = 0;

    #pragma omp parallel
    {
        #pragma omp single
        {
            #pragma omp task
            i = 1;
            #pragma omp task
            i = 2;
        }
    }

    printf("i = %d\n", i);

    return 0;
} 