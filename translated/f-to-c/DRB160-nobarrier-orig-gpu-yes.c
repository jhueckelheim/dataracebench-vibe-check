/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Vector addition followed by multiplication involving the same var should have a barrier in between.
//omp distribute directive does not have implicit barrier. This will cause data race.
//Data Race Pair: b[i]@36:23:R vs. b[i]@42:13:W

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
int a, i, j, k, val;
int b[8], c[8], temp[8];

int main()
{
    for (i = 0; i < 8; i++) {
        b[i] = 0;
        c[i] = 2;
        temp[i] = 0;
    }

    a = 2;
    val = 0;

    #pragma omp target map(tofrom:b[0:8]) map(to:c[0:8],temp[0:8],a) device(0)
    #pragma omp teams
    for (i = 0; i < 100; i++) {
        #pragma omp distribute
        for (j = 0; j < 8; j++) {
            temp[j] = b[j] + c[j];
        }
        // No implicit barrier after distribute

        #pragma omp distribute
        for (j = 7; j >= 0; j--) {  // Adjust for 0-based indexing
            b[j] = temp[j] * a;  // Data race: temp may not be updated yet
        }
    }

    for (i = 0; i < 100; i++) {
        val = val + 2;
        val = val * 2;
    }

    for (i = 0; i < 8; i++) {
        if (val != b[i]) {
            printf("%d %d\n", b[i], val);
        }
    }

    return 0;
} 