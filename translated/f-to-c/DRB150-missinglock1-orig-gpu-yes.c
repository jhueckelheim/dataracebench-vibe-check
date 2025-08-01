/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The distribute parallel do directive at line 22 will execute loop using multiple teams.
//The loop iterations are distributed across the teams in chunks in round robin fashion.
//The omp lock is only guaranteed for a contention group, i.e, within a team.
//Data Race Pair, var@25:9:W vs. var@25:9:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i;
    omp_lock_t lck;
    omp_init_lock(&lck);

    #pragma omp target map(tofrom:var) device(0)
    #pragma omp teams distribute parallel for
    for (i = 0; i < 10; i++) {
        omp_set_lock(&lck);
        var = var + 1;  // Data race: lock only works within team, not across teams
        omp_unset_lock(&lck);
    }

    omp_destroy_lock(&lck);

    printf("%d\n", var);

    return 0;
} 