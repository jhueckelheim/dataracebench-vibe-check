/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Race condition on outLen due to unprotected writes.
//Adding private (outLen) can avoid race condition. But it is wrong semantically.
//
//Data race pairs: we allow two pair to preserve the original code pattern.
//1. outLen@60:9:W vs. outLen@60:9:W
//2. output[]@59:9:W vs. output[]@59:9:W

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int i, inLen, outLen;
    int* input;
    int* output;

    inLen = 1000;
    outLen = 0;  // Adjust for 0-based indexing

    if (argc == 1) {
        printf("No command line arguments provided.\n");
    }

    if (argc >= 2) {
        inLen = atoi(argv[1]);
        if (inLen <= 0) {
            printf("Error, invalid integer value.\n");
        }
    }

    input = (int*)malloc(inLen * sizeof(int));
    output = (int*)malloc(inLen * sizeof(int));

    for (i = 0; i < inLen; i++) {
        input[i] = i + 1;
    }

    #pragma omp parallel for
    for (i = 0; i < inLen; i++) {
        output[outLen] = input[i];
        outLen = outLen + 1;
    }

    printf("output(0) = %d\n", output[0]);

    free(input);
    free(output);

    return 0;
} 