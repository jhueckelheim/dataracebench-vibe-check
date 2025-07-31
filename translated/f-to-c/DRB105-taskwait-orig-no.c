/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

/* This is a program based on a test contributed by Yizi Gu@Rice Univ.
 * Classic Fibonacci calculation using task+taskwait. No data races.
 */

#include <omp.h>
#include <stdio.h>

// Global variable (from module)
int input;

int fib(int n)
{
    int i, j, r;

    if (n < 2) {
        r = n;
    } else {
        #pragma omp task shared(i)
        i = fib(n - 1);

        #pragma omp task shared(j)
        j = fib(n - 2);

        #pragma omp taskwait
        r = i + j;
    }
    return r;
}

int main()
{
    int result;
    input = 30;

    #pragma omp parallel
    {
        #pragma omp single
        result = fib(input);
    }

    printf("Fib for %d = %d\n", input, result);

    return 0;
} 