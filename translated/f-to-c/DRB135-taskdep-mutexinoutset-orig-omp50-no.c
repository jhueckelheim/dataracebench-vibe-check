/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Addition of mutexinoutset dependence type on c, will ensure that line d@36:9 assignment will depend
//on task at Line 29 and line 32. They might execute in any order but not at the same time.
//There is no data race.

#include <omp.h>
#include <stdio.h>

int main()
{
    int a, b, c, d;

    #pragma omp parallel
    {
        #pragma omp single
        {
            #pragma omp task depend(out: c)
            c = 1;  // Task T1

            #pragma omp task depend(out: a)
            a = 2;  // Task T2

            #pragma omp task depend(out: b)
            b = 3;  // Task T3

            #pragma omp task depend(in: a) depend(mutexinoutset: c)
            c = c + a;  // Task T4

            #pragma omp task depend(in: b) depend(mutexinoutset: c)
            c = c + b;  // Task T5

            #pragma omp task depend(in: c)
            d = c;  // Task T6
        }
    }

    printf("%d\n", d);

    return 0;
} 