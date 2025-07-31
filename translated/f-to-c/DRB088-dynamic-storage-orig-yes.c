/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//For the case of a variable which is not referenced within a construct:
//objects with dynamic storage duration should be shared.
//Putting it within a threadprivate directive may cause seg fault since
// threadprivate copies are not allocated!
//
//Dependence pair: *counter@22:9:W vs. *counter@22:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global pointer (from module)
int* counter;

void foo()
{
    *counter = *counter + 1;  // Data race on dynamically allocated memory
}

int main()
{
    counter = (int*)malloc(sizeof(int));
    *counter = 0;

    #pragma omp parallel
    foo();

    printf("%d\n", *counter);

    free(counter);

    return 0;
} 