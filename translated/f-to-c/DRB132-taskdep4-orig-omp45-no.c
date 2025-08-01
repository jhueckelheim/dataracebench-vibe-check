/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The second taskwait ensures that the second child task has completed; hence it is safe to access
//the y variable in the following print statement.

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
    y = y - 1;  // 2nd child task

    #pragma omp task depend(in: x) if(0)  // 1st taskwait (undeferred)
    {
        // Empty task body - acts as taskwait for x-dependent tasks only
    }

    printf("x = %d\n", x);

    #pragma omp taskwait  // 2nd taskwait - ensures all tasks complete

    printf("y = %d\n", y);  // No race: taskwait ensures completion
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