/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Loop carried true dep between tmp =..  and ..= tmp.
//Data race pair: tmp@48:16:R vs. tmp@49:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, tmp, len;
    int* a;

    len = 100;
    tmp = 10;

    if (argc == 1) {
        printf("No command line arguments provided.\n");
    }

    if (argc >= 2) {
        len = atoi(argv[1]);
        if (len <= 0) {
            printf("Error, invalid integer value.\n");
        }
    }

    a = (int*)malloc(len * sizeof(int));

    #pragma omp parallel for
    for (i = 0; i < len; i++) {
        a[i] = tmp;
        tmp = a[i] + i + 1;  // Adjust for 0-based indexing
    }

    free(a);

    return 0;
} 