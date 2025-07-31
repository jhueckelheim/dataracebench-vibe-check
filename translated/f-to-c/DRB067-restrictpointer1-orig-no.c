/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Array initialization using assignments. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

void foo(double** newSxx, double** newSyy, int len)
{
    int i;
    double* tar1;
    double* tar2;

    *newSxx = (double*)malloc(len * sizeof(double));
    *newSyy = (double*)malloc(len * sizeof(double));

    tar1 = (double*)malloc(len * sizeof(double));
    tar2 = (double*)malloc(len * sizeof(double));

    *newSxx = tar1;
    *newSyy = tar2;

    #pragma omp parallel for private(i) firstprivate(len)
    for (i = 0; i < len; i++) {
        tar1[i] = 0.0;
        tar2[i] = 0.0;
    }

    printf("%f %f\n", tar1[len-1], tar2[len-1]);
    
    // In C, newSxx and newSyy are pointing to tar1 and tar2,
    // so we only need to free tar1 and tar2
    free(tar1);
    free(tar2);
}

int main()
{
    int len = 1000;
    double* newSxx;
    double* newSyy;

    newSxx = (double*)malloc(len * sizeof(double));
    newSyy = (double*)malloc(len * sizeof(double));

    foo(&newSxx, &newSyy, len);

    // newSxx and newSyy now point to memory allocated inside foo(),
    // which has already been freed
    return 0;
} 