/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

// Iteration 1 and 2 can have conflicting writes to a(1). But if they are scheduled to be run by 
// the same thread, dynamic tools may miss this.
// Data Race Pair, a(0)@39:9:W vs. a(i)@40:22:W

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

int main()
{
    int i, N;
    int* a;

    N = 100;

    a = (int*)malloc(N * sizeof(int));

    load_from_input(a, N);

    #pragma omp parallel for shared(a)
    for (i = 0; i < N; i++) {
        a[i] = i + 1;  // Adjust for 0-based indexing
        if (i == 1) a[0] = 1;  // Data race: thread writing a[1] also writes a[0]
    }

    free(a);

    return 0;
} 