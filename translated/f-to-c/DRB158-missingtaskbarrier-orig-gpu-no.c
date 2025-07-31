/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Depend clause at line 29 and 33 will ensure that there is no data race.

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
int a, i;
int x[64], y[64];

int main()
{
    for (i = 0; i < 64; i++) {
        x[i] = 0;
        y[i] = 3;
    }

    a = 5;

    #pragma omp target map(to:y[0:64],a) map(tofrom:x[0:64]) device(0)
    for (i = 0; i < 64; i++) {
        #pragma omp task depend(inout:x[i])
        x[i] = a * x[i];

        #pragma omp task depend(inout:x[i])  // Dependency ensures no race
        x[i] = x[i] + y[i];
    }

    for (i = 0; i < 64; i++) {
        if (x[i] != 3) {
            printf("%d\n", x[i]);
        }
    }

    #pragma omp taskwait

    return 0;
} 