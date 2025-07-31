/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two tasks with depend clause to ensure execution order:
//i is shared for two tasks based on implicit data-sharing attribute rules. No data race pairs.

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
            #pragma omp task depend(out:i)
            i = 1;

            #pragma omp task depend(in:i)
            i = 2;
        }
    }

    if (i != 2) {
        printf("i is not equal to 2\n");
    }

    return 0;
} 