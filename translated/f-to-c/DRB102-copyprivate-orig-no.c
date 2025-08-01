/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//threadprivate+copyprivate: no data races

#include <omp.h>
#include <stdio.h>

// Global variables (from module)
int y;
double x;

#pragma omp threadprivate(x, y)

int main()
{
    #pragma omp parallel
    {
        #pragma omp single copyprivate(x, y)
        {
            x = 1.0;
            y = 1;
        }
    }

    printf("x = %.1f  y = %d\n", x, y);

    return 0;
} 