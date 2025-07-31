/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A master directive is used to protect memory accesses. No data race pairs.

#include <omp.h>
#include <stdio.h>

int main()
{
    int k;

    #pragma omp parallel
    {
        #pragma omp master
        {
            k = omp_get_num_threads();
            printf("Number of threads requested = %d\n", k);
        }
    }

    return 0;
} 