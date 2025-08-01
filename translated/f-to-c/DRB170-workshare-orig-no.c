/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The workshare construct is only available in Fortran. The workshare spreads work across the threads 
//executing the parallel. There is an implicit barrier. No data race.

#include <omp.h>
#include <stdio.h>

int main()
{
    int AA, BB, CC, res;

    BB = 1;
    CC = 2;

    #pragma omp parallel
    {
        #pragma omp single  // Simulating workshare with implicit barrier
        {
            AA = BB;
            AA = AA + CC;
        }
        // Implicit barrier here

        #pragma omp single  // Simulating second workshare
        {
            res = AA * 2;  // No race: barrier ensures AA is updated
        }
    }

    if (res != 6) {
        printf("%d\n", res);
    }

    return 0;
} 