/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//No data race. The data environment of the task is created according to the
//data-sharing attribute clauses, here at line 21:27 it is var. Hence, var is
//modified 10 times, resulting to the value 10.

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
                #pragma omp task shared(var) if(0)  // if(0) makes tasks undeferred
                var = var + 1;
            }
        }
    }

    printf("var = %d\n", var);

    return 0;
} 