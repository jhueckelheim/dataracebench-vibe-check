/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Concurrent access on same variable var@23 and var@25 leads to the race condition if two different
//locks are used. This is the reason here we have used the atomic directive to ensure that addition
//and subtraction are not interleaved. No data race pairs.

#include <omp.h>
#include <stdio.h>

int main()
{
    int var, i;
    var = 0;

    #pragma omp target map(tofrom:var) device(0)
    #pragma omp teams distribute parallel for
    for (i = 0; i <= 100; i++) {
        #pragma omp atomic
        var = var + 1;
        #pragma omp atomic
        var = var - 2;  // Atomic operations prevent interleaving
    }

    printf("%d\n", var);

    return 0;
} 