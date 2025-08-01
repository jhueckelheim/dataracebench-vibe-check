/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Classic PI calculation using reduction. No data race pairs.

#include <omp.h>
#include <stdio.h>

int main()
{
    long double x, interval_width, pi;
    long long i, num_steps;

    pi = 0.0;
    num_steps = 2000000000LL;
    interval_width = 1.0 / num_steps;

    #pragma omp parallel for reduction(+:pi) private(x)
    for (i = 0; i < num_steps; i++) {
        x = (i + 1 + 0.5) * interval_width;  // Adjust for 0-based indexing
        pi = pi + 1.0 / (x*x + 1.0);
    }

    pi = pi * 4.0 * interval_width;
    printf("PI = %.20Lf\n", pi);

    return 0;
} 