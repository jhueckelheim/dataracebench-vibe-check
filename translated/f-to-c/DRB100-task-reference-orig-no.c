/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

/* Cover the implicitly determined rule: In an orphaned task generating construct,
 * formal arguments passed by reference are firstprivate.
 * This requires OpenMP 4.5 to work.
 * Earlier OpenMP does not allow a reference type for a variable within firstprivate().
 * No data race pairs.
 */

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global variable (from module)
int* a;

void gen_task(int* i_ptr)  // Passed by reference
{
    #pragma omp task firstprivate(i_ptr)  // Reference parameters are firstprivate
    {
        int i = *i_ptr;  // Dereference the pointer
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
                gen_task(&i);  // Pass by reference
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