/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//tasks with depend clauses to ensure execution order, no data races.

#include <omp.h>
#include <stdio.h>
#include <unistd.h>

int main()
{
    int i, j, k;
    i = 0;

    #pragma omp parallel
    {
        #pragma omp single
        {
            #pragma omp task depend(out:i)
            {
                sleep(3);
                i = 1;
            }
            
            #pragma omp task depend(in:i)
            {
                j = i;
            }
            
            #pragma omp task depend(in:i)
            {
                k = i;
            }
        }
    }

    printf("j = %d  k = %d\n", j, k);

    if (j != 1 && k != 1) {
        printf("Race Condition\n");
    }

    return 0;
} 