/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Concurrent access of var@22:5 in an intra region. Missing Lock leads to intra region data race.
//Data Race pairs, var@22:13:W vs. var@22:13:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i;
    var = 0;

    #pragma omp target map(tofrom:var) device(0)
    #pragma omp teams num_teams(1)
    #pragma omp distribute parallel for
    for (i = 0; i < 100; i++) {
        var = var + 1;  // Data race: no lock synchronization
    }

    printf("%d\n", var);

    return 0;
} 