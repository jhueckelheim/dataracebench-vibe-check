/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//omp_get_thread_num() is used to ensure serial semantics. No data race pairs.

#include <omp.h>
#include <stdio.h>

int main()
{
    int numThreads;

    #pragma omp parallel
    {
        if (omp_get_thread_num() == 0) {
            numThreads = omp_get_num_threads();
        }
    }

    printf("numThreads = %d\n", numThreads);

    return 0;
} 