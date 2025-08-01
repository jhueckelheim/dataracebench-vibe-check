/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Test if the semantics of omp_get_thread_num() is correctly recognized.
//Thread with id 0 writes numThreads while other threads read it, causing data races.
//Data race pair: numThreads@22:9:W vs. numThreads@24:31:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int numThreads;
    numThreads = 0;

    #pragma omp parallel
    {
        if (omp_get_thread_num() == 0) {
            numThreads = omp_get_num_threads();
        } else {
            printf("numThreads = %d\n", numThreads);
        }
    }

    return 0;
} 