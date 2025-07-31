/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//For a variable declared in a scope inside an OpenMP construct:
//* private if the variable has an automatic storage duration
//* shared if the variable has a static storage duration.
//
//Dependence pairs:
//   tmp@38:13:W vs. tmp@38:13:W
//   tmp@38:13:W vs. tmp@39:20:R

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

static int tmp = 0;  // Static storage duration - shared, causes data race

int main()
{
    int i, len;
    int* a;
    int* b;
    int tmp2;  // Automatic storage duration - private, no race

    len = 100;
    a = (int*)malloc(len * sizeof(int));
    b = (int*)malloc(len * sizeof(int));

    for (i = 0; i < len; i++) {
        a[i] = i + 1;  // Adjust for 0-based indexing
        b[i] = i + 1;
    }

    #pragma omp parallel
    {
        #pragma omp for
        for (i = 0; i < len; i++) {
            tmp = a[i] + (i + 1);  // Data race - static storage
            a[i] = tmp;
        }
    }

    #pragma omp parallel
    {
        #pragma omp for
        for (i = 0; i < len; i++) {
            tmp2 = b[i] + (i + 1);  // No race - automatic storage
            b[i] = tmp2;
        }
    }

    printf("%d   %d\n", a[49], b[49]);  // Adjust for 0-based indexing

    free(a);
    free(b);

    return 0;
} 