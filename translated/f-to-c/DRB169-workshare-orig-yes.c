/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The workshare construct is only available in Fortran. The workshare spreads work across the threads 
//executing the parallel. There is an implicit barrier. The nowait nullifies this barrier and hence
//there is a race at line:29 due to nowait at line:26. Data Race Pairs, AA@25:9:W vs. AA@29:15:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int AA, BB, CC, res;

    BB = 1;
    CC = 2;

    #pragma omp parallel
    {
        #pragma omp single nowait  // Simulating workshare nowait
        {
            AA = BB;
            AA = AA + CC;
        }
        // No barrier due to nowait

        #pragma omp single  // Simulating second workshare
        {
            res = AA * 2;  // Data race: AA may not be updated yet
        }
    }

    if (res != 6) {
        printf("%d\n", res);
    }

    return 0;
} 