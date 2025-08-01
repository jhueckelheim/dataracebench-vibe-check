/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

/* This is a program based on a test contributed by Yizi Gu@Rice Univ.
 * Use taskgroup to synchronize two tasks. No data race pairs.
 */

#include <omp.h>
#include <stdio.h>
#include <unistd.h>

int main()
{
    int result;
    result = 0;

    #pragma omp parallel
    {
        #pragma omp single
        {
            #pragma omp taskgroup
            {
                #pragma omp task
                {
                    sleep(3);
                    result = 1;
                }
            }
            
            #pragma omp task
            result = 2;
        }
    }

    printf("result = %d\n", result);

    return 0;
} 