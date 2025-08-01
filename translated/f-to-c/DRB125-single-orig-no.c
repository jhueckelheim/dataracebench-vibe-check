/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This example is derived from an example by Simone Atzeni, NVIDIA.
//
//Description: Race on variable init if used master construct. The variable is written by the
//master thread and concurrently read by the others.
//
//Solution: master construct does not have an implicit barrier better
//use single at line 26. Fixed version for DRB124-master-orig-yes.c. No data race.

#include <omp.h>
#include <stdio.h>

int main()
{
    int init, local;

    #pragma omp parallel shared(init) private(local)
    {
        #pragma omp single
        init = 10;  // Single has implicit barrier, ensuring synchronization

        local = init;  // No race: all threads see init=10 after the barrier
    }

    return 0;
} 