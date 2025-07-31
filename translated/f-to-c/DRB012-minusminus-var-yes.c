/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//The -1 operation is not protected, causing race condition.
//Data race pair: numNodes2@59:13:W vs. numNodes2@59:13:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, len, numNodes, numNodes2;
    int* x;

    len = 100;

    if (argc == 1) {
        printf("No command line arguments provided.\n");
    }

    if (argc >= 2) {
        len = atoi(argv[1]);
        if (len <= 0) {
            printf("Error, invalid integer value.\n");
        }
    }

    x = (int*)malloc(len * sizeof(int));

    numNodes = len;
    numNodes2 = 0;

    // initialize x()
    for (i = 0; i < len; i++) {
        if ((i+1) % 2 == 0) {
            x[i] = 5;
        } else {
            x[i] = -5;
        }
    }

    #pragma omp parallel for
    for (i = numNodes-1; i >= 0; i--) {
        if (x[i] <= 0) {
            numNodes2 = numNodes2 - 1;
        }
    }

    printf("numNodes2 = %d\n", numNodes2);

    free(x);

    return 0;
} 