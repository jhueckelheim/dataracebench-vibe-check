/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

// Input dependence race: example from OMPRacer: A Scalable and Precise Static Race
// Detector for OpenMP Programs
// Data Race Pair, a(1)@63:26:W vs. a(i)@62:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

void load_from_input(int* a, int N)
{
    int i;
    for (i = 0; i < N; i++) {
        a[i] = i + 1;  // Adjust for 0-based indexing
    }
}

int main(int argc, char* argv[])
{
    int i, N;
    int* a;

    N = 100;

    if (argc == 0) {
        printf("No command line arguments provided.\n");
    }

    if (argc >= 2) {
        N = atoi(argv[1]);
        if (N <= 0) {
            printf("Error, invalid integer value.\n");
        }
    }

    a = (int*)malloc(N * sizeof(int));

    load_from_input(a, N);

    #pragma omp parallel for shared(a)
    for (i = 0; i < N; i++) {
        a[i] = i + 1;  // Adjust for 0-based indexing
        if (N > 10000) a[0] = 1;  // Data race: multiple threads write to a[0]
    }

    free(a);

    return 0;
} 