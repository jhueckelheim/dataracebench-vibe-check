/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//use of omp target + teams + distribute + parallel for. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    long long i, i2, len, l_limit, tmp;
    double sum, sum2;
    double* a;
    double* b;

    len = 2560;
    sum = 0.0;
    sum2 = 0.0;

    a = (double*)malloc(len * sizeof(double));
    b = (double*)malloc(len * sizeof(double));

    for (i = 0; i < len; i++) {
        a[i] = (double)(i + 1) / 2.0;  // Adjust for 0-based indexing
        b[i] = (double)(i + 1) / 3.0;
    }

    #pragma omp target map(to: a[0:len], b[0:len]) map(tofrom: sum)
    #pragma omp teams num_teams(10) thread_limit(256) reduction(+:sum)
    #pragma omp distribute
    for (i2 = 0; i2 < len; i2 += 256) {
        #pragma omp parallel for reduction(+:sum)
        for (i = i2 + 1; i <= (i2 + 256 < len ? i2 + 256 : len); i++) {  // Adjust bounds
            sum = sum + a[i-1] * b[i-1];  // Adjust for 0-based indexing
        }
    }

    #pragma omp parallel for reduction(+:sum2)
    for (i = 0; i < len; i++) {
        sum2 = sum2 + a[i] * b[i];
    }

    printf("sum = %d; sum2 = %d\n", (int)sum, (int)sum2);

    free(a);
    free(b);

    return 0;
} 