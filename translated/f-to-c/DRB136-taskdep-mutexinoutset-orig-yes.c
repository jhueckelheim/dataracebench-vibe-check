/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Due to the missing mutexinoutset dependence type on c, these tasks will execute in any
//order leading to the data race at line 35. Data Race Pair, d@35:9:W vs. d@35:9:W

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

            #pragma omp task depend(in: a)  // Missing mutexinoutset for c
            c = c + a;  // Task T4 - data race on c

            #pragma omp task depend(in: b)  // Missing mutexinoutset for c
            c = c + b;  // Task T5 - data race on c

            #pragma omp task depend(in: c)
            d = c;  // Task T6
        }
    }

    printf("%d\n", d);

    return 0;
} 