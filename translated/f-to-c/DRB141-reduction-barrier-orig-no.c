/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//To avoid data race, the initialization of the original list item "a" should complete before any
//update of a as a result of the reduction clause. This can be achieved by adding an explicit
//barrier after the assignment a=0@22:9, or by enclosing the assignment a=0@22:9 in a single directive
//or by initializing a@21:7 before the start of the parallel region. No data race pair

#include <omp.h>
#include <stdio.h>

int main()
{
    int a, i;

    #pragma omp parallel shared(a) private(i)
    {
        #pragma omp master
        a = 0;

        #pragma omp barrier  // Explicit barrier fixes the race

        #pragma omp for reduction(+:a)
        for (i = 1; i <= 10; i++) {
            a = a + i;
        }

        #pragma omp single
        printf("Sum is %d\n", a);
    }

    return 0;
} 