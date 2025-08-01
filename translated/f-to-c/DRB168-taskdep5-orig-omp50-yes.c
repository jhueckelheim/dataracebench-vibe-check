/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//There is no completion restraint on the second child task. Hence, immediately after the first
//taskwait it is unsafe to access the y variable since the second child task may still be
//executing.
//Data Race at y@34:8:W vs. y@40:23:R

#include <omp.h>
#include <stdio.h>

void foo()
{
    int x, y;
    x = 0;
    y = 2;

    #pragma omp task depend(inout: x) shared(x)
    x = x + 1;  // 1st Child Task

    #pragma omp task shared(y)
    y = y - x;  // 2nd child task

    #pragma omp taskwait depend(in: x)  // 1st taskwait - OpenMP 5.0 feature

    printf("x = %d\n", x);
    printf("y = %d\n", y);  // Data race: 2nd task may still be running

    #pragma omp taskwait  // 2nd taskwait
}

int main()
{
    #pragma omp parallel
    {
        #pragma omp single
        foo();
    }

    return 0;
} 