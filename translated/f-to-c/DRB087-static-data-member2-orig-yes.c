/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//For the case of a variable which is referenced within a construct:
//static data member should be shared, unless it is within a threadprivate directive.
//
//Dependence pair: counter@37:5:W vs. counter@37:5:W

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
int counter = 0;
int pcounter = 0;

#pragma omp threadprivate(pcounter)

typedef struct {
    int counter;
    int pcounter;
} A;

int main()
{
    A c = {0, 0};

    #pragma omp parallel
    {
        counter = counter + 1;    // Data race - not threadprivate
        pcounter = pcounter + 1;  // No race - threadprivate
    }

    printf("%d %d\n", counter, pcounter);

    return 0;
} 