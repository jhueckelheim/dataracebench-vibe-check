/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two tasks with a lock synchronization to ensure execution order. No data race pairs.

#include <omp.h>
#include <stdio.h>

int main()
{
    omp_lock_t lock;
    int i;

    i = 0;
    omp_init_lock(&lock);

    #pragma omp parallel sections
    {
        #pragma omp section
        {
            omp_set_lock(&lock);
            i = i + 1;
            omp_unset_lock(&lock);
        }
        #pragma omp section
        {
            omp_set_lock(&lock);
            i = i + 2;
            omp_unset_lock(&lock);
        }
    }

    omp_destroy_lock(&lock);

    printf("I = %d\n", i);

    return 0;
} 