/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The -1 operation on numNodes2 is not protected, causing data race.
//Data race pair: numNodes2@32:13:W vs. numNodes2@32:13:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, len, numNodes, numNodes2;
    int x[100];
    len = 100;
    numNodes = len;
    numNodes2 = 0;

    for (i = 0; i < len; i++) {
        if ((i+1) % 2 == 0) {
            x[i] = 5;
        } else {
            x[i] = -5;
        }
    }

    #pragma omp parallel for
    for (i = numNodes-1; i >= 0; i--) {
        if (x[i] <= 0) {
            numNodes2 = numNodes2 - 1;
        }
    }

    printf("numNodes2 = %d\n", numNodes2);

    return 0;
} 