/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The distribute parallel do directive at line 19 will execute loop using multiple teams.
//The loop iterations are distributed across the teams in chunks in round robin fashion.
//The missing lock enclosing var@21 leads to data race. Data Race Pairs, var@21:9:W vs. var@21:9:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i;

    #pragma omp target map(tofrom:var) device(0)
    #pragma omp teams distribute parallel for
    for (i = 0; i < 100; i++) {
        var = var + 1;  // Data race: no synchronization across teams
    }

    printf("%d\n", var);

    return 0;
} 