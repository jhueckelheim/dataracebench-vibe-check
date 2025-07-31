/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two tasks without synchronization to protect data write, causing data races.
//Data race pair: i@20:5:W vs. i@22:5:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int i;
    i = 0;

    #pragma omp parallel sections
    {
        #pragma omp section
        i = 1;
        #pragma omp section
        i = 2;
    }

    printf("i = %d\n", i);

    return 0;
} 