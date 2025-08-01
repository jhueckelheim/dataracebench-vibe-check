/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Using lastprivate() to resolve an output dependence.
//
//Semantics of lastprivate (x):
//causes the corresponding original list item to be updated after the end of the region.
//The compiler/runtime copies the local value back to the shared one within the last iteration.

#include <omp.h>
#include <stdio.h>

void foo()
{
    int i, x;
    
    #pragma omp parallel for private(i) lastprivate(x)
    for (i = 0; i < 100; i++) {
        x = i + 1;  // Adjust for 1-based indexing in output
    }
    
    printf("x = %d\n", x);
}

int main()
{
    #pragma omp parallel
    foo();

    return 0;
} 