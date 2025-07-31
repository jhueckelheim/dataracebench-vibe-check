/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Though we have used critical directive to ensure that additions across teams are not overlapped.
//Critical only synchronizes within a team. There is a data race pair.
//Data Race pairs, var@24:9:W vs. var@24:15:R

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i;
    var = 0;

    #pragma omp target map(tofrom:var) device(0)
    #pragma omp teams distribute parallel for
    for (i = 0; i < 100; i++) {
        #pragma omp critical(addlock)
        var = var + 1;  // Data race: critical only within team, not across teams
    }

    printf("%d\n", var);

    return 0;
} 