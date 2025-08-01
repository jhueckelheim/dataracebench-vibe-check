/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Concurrent access of var@23 in an intra region. Lock ensures that there is no data race.

#include <omp.h>
#include <stdio.h>

int main()
{
    omp_lock_t lck;
    int var, i;
    var = 0;

    #pragma omp target map(tofrom:var) device(0)
    #pragma omp teams num_teams(1)
    #pragma omp distribute parallel for
    for (i = 0; i < 100; i++) {
        omp_set_lock(&lck);
        var = var + 1;  // No race: lock synchronizes within team
        omp_unset_lock(&lck);
    }

    return 0;
} 