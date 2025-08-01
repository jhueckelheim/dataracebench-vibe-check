/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The increment at line number 22 is critical for the variable
//var@22:13. Therefore, there is a possible Data Race pair var@22:13:W vs. var@22:19:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i;
    var = 0;

    #pragma omp target map(tofrom:var) device(0)
    #pragma omp teams distribute parallel for
    for (i = 0; i < 200; i++) {
        #pragma omp critical
        var = var + 1;  // Data race across teams - missing reduction
    }

    printf("%d\n", var);

    return 0;
} 