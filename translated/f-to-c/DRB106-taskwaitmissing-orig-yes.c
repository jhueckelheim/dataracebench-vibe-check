/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

/* This is a program based on a test contributed by Yizi Gu@Rice Univ.
 * Classic Fibonacci calculation using task but missing taskwait.
 * Data races pairs: i@29:13:W vs. i@34:17:R
 *                   j@32:13:W vs. j@34:19:R
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

        r = i + j;  // Data race: missing taskwait before this line
    }
    
    #pragma omp taskwait  // Misplaced - should be before r = i + j
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