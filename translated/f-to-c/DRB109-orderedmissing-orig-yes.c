/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This is a program based on a test contributed by Yizi Gu@Rice Univ.
/* Missing the ordered clause
 * Data race pair: x@21:9:W vs. x@21:9:W
 */

#include <omp.h>
#include <stdio.h>

int main()
{
    int x, i;
    x = 0;

    #pragma omp parallel for ordered
    for (i = 0; i < 100; i++) {
        // Missing ordered directive here causes data race
        x = x + 1;
    }

    printf("x = %d\n", x);

    return 0;
} 