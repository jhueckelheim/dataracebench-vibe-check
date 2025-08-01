/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This is a program based on a test contributed by Yizi Gu@Rice Univ.
//Proper user of ordered directive and clause, no data races

#include <omp.h>
#include <stdio.h>

int main()
{
    int x, i;
    x = 0;

    #pragma omp parallel for ordered
    for (i = 0; i < 100; i++) {
        #pragma omp ordered
        x = x + 1;
    }

    printf("x = %d\n", x);

    return 0;
} 