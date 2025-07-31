/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The assignment to a@21:9 is  not synchronized with the update of a@29:11 as a result of the
//reduction computation in the for loop.
//Data Race pair: a@21:9:W vs. a@24:30:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int a, i;

    #pragma omp parallel shared(a) private(i)
    {
        #pragma omp master
        a = 0;  // Data race: no barrier before reduction

        #pragma omp for reduction(+:a)
        for (i = 1; i <= 10; i++) {
            a = a + i;
        }

        #pragma omp single
        printf("Sum is %d\n", a);
    }

    return 0;
} 