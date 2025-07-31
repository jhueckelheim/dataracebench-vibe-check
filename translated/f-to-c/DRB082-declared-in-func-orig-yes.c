/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A variable is declared inside a function called within a parallel region.
//The variable should be shared if it uses static storage.
//
//Data race pair: i@19:7:W vs. i@19:7:W

#include <omp.h>
#include <stdio.h>

void foo()
{
    static int i = 0;  // Static storage, shared among threads
    i = i + 1;
    printf("%d\n", i);
}

int main()
{
    #pragma omp parallel
    foo();

    return 0;
} 