/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A variable is declared inside a function called within a parallel region.
//The variable should be private if it does not use static storage. No data race pairs.

#include <omp.h>
#include <stdio.h>

void foo()
{
    int q;  // Automatic storage, private to each thread
    q = 0;
    q = q + 1;
}

int main()
{
    #pragma omp parallel
    foo();

    return 0;
} 