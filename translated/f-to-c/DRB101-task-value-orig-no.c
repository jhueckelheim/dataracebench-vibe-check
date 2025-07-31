/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

/* Cover an implicitly determined rule: In a task generating construct,
 * a variable without applicable rules is firstprivate. No data race pairs.
 */

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global variable (from module)
int* a;

void gen_task(int i)  // Passed by value
{
    #pragma omp task  // i is implicitly firstprivate when passed by value
    {
        a[i-1] = i + 1;  // Adjust for 0-based indexing
    }
}

int main()
{
    int i;
    
    a = (int*)malloc(100 * sizeof(int));

    #pragma omp parallel
    {
        #pragma omp single
        {
            for (i = 1; i <= 100; i++) {
                gen_task(i);  // Pass by value
            }
        }
    }

    for (i = 1; i <= 100; i++) {
        if (a[i-1] != i + 1) {  // Adjust for 0-based indexing
            printf("warning: a(%d) = %d not expected %d\n", i, a[i-1], i + 1);
        }
    }

    free(a);

    return 0;
} 