/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Freshly allocated pointers do not alias to each other. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

void setup(int N)
{
    int i;
    double* m_pdv_sum;
    double* m_nvol;
    double* tar1;
    double* tar2;

    m_pdv_sum = (double*)malloc(N * sizeof(double));
    m_nvol = (double*)malloc(N * sizeof(double));
    tar1 = (double*)malloc(N * sizeof(double));
    tar2 = (double*)malloc(N * sizeof(double));

    m_pdv_sum = tar1;
    m_nvol = tar2;

    #pragma omp parallel for schedule(static)
    for (i = 0; i < N; i++) {
        tar1[i] = 0.0;
        tar2[i] = (i + 1) * 2.5;  // Adjust for 0-based indexing
    }

    // printf("%f %f\n", tar1[N-1], tar2[N-1]);
    
    // In C, we just need to free the allocated memory
    free(tar1);
    free(tar2);
    // Note: m_pdv_sum and m_nvol were pointing to tar1 and tar2, 
    // so we don't free them separately
}

int main()
{
    int N = 1000;
    setup(N);
    return 0;
} 