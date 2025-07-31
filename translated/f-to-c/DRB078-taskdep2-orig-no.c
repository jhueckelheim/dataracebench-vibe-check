/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two tasks with depend clause to ensure execution order, no data races.
//i is shared for two tasks based on implicit data-sharing attribute rules.

#include <omp.h>
#include <stdio.h>
#include <unistd.h>

int main()
{
    int i;
    i = 0;

    #pragma omp parallel
    {
        #pragma omp single
        {
            #pragma omp task depend(out:i)
            {
                sleep(3);
                i = 3;
            }
            
            #pragma omp task depend(out:i)
            {
                i = 2;
            }
        }
    }

    if (i != 2) {
        printf("%d\n", i);
    }

    return 0;
} 