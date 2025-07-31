/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The thread encountering the taskwait directive at line 22 only waits for its child task
//(line 14-21) to complete. It does not wait for its descendant tasks (line 16-19). Data Race pairs, sum@36:13:W vs. sum@36:13:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int* a;
    int* psum;
    int sum, i;

    a = (int*)malloc(4 * sizeof(int));
    psum = (int*)malloc(4 * sizeof(int));

    #pragma omp parallel num_threads(2)
    {
        #pragma omp for schedule(dynamic, 1)
        for (i = 0; i < 4; i++) {
            a[i] = i + 1;  // Adjust for 0-based indexing
        }

        #pragma omp single
        {
            #pragma omp task
            {
                #pragma omp task
                psum[1] = a[2] + a[3];  // Descendant task

                psum[0] = a[0] + a[1];  // Child task
            }
            #pragma omp taskwait  // Only waits for child, not descendant
            sum = psum[1] + psum[0];  // Data race if descendant not complete
        }
    }

    printf("sum = %d\n", sum);

    free(a);
    free(psum);

    return 0;
} 