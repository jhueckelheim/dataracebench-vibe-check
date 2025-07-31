/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A single thread will spawn all the tasks. Add if(0) to avoid the data race, undeferring the tasks.
//Data Race Pairs, var@21:9:W vs. var@21:9:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i;
    var = 0;

    #pragma omp parallel sections
    {
        #pragma omp section
        {
            for (i = 1; i <= 10; i++) {
                #pragma omp task shared(var)  // Missing if(0) - tasks are deferred
                var = var + 1;  // Data race among deferred tasks
            }
        }
    }

    printf("var = %d\n", var);

    return 0;
} 