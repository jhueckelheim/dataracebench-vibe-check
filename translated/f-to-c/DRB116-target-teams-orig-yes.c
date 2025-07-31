/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//use of omp target + teams
//Without protection, master threads from two teams cause data races.
//Data race pair: a@24:9:W vs. a@24:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int i, len;
    double* a;

    len = 100;
    a = (double*)malloc(len * sizeof(double));

    for (i = 0; i < len; i++) {
        a[i] = (double)(i + 1) / 2.0;  // Adjust for 0-based indexing
    }

    #pragma omp target map(tofrom: a[0:len])
    #pragma omp teams num_teams(2)
    a[49] = a[49] * 2.0;  // Data race - both teams write to same location

    printf("a(50) = %f\n", a[49]);  // Adjust for 0-based indexing

    free(a);

    return 0;
} 