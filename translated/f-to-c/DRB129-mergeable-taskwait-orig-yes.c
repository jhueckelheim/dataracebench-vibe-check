/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Taken from OpenMP Examples 5.0, example tasking.12.c
//The created task will access different instances of the variable x if the task is not merged,
//as x is firstprivate, but it will access the same variable x if the task is merged. It can
//Data Race Pairs, x@22:5:W vs. x@22:5:W
//print two different values for x depending on the decisions taken by the implementation.

#include <omp.h>
#include <stdio.h>

int main()
{
    int x;
    x = 2;

    #pragma omp task mergeable  // x is implicitly firstprivate, causing race on merge
    x = x + 1;

    printf("x = %d\n", x);

    return 0;
} 