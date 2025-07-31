/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//A function argument passed by value should be private inside the function.
//Variable i is read only. No data race pairs.

#include <omp.h>
#include <stdio.h>

void f1(int i)
{
    i = i + 1;
}

int main()
{
    int i;
    i = 0;

    #pragma omp parallel
    f1(i);

    if (i != 0) {
        printf("i = %d\n", i);
    }

    return 0;
} 