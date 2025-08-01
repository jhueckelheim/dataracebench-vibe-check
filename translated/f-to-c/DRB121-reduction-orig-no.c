/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Number of threads is empirical: We need enough threads so that
//the reduction is really performed hierarchically in the barrier!
//There is no data race.

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i, sum1, sum2;

    var = 0;
    sum1 = 0;
    sum2 = 0;

    #pragma omp parallel reduction(+: var)
    {
        #pragma omp for schedule(static) reduction(+: sum1)
        for (i = 1; i <= 5; i++) {
            sum1 = sum1 + i;
        }

        #pragma omp for schedule(static) reduction(+: sum2)
        for (i = 1; i <= 5; i++) {
            sum2 = sum2 + i;
        }

        var = sum1 + sum2;
    }

    printf("var = %d\n", var);

    return 0;
} 