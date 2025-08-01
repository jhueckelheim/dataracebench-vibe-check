/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Data race on outLen due to ++ operation.
//Adding private (outLen) can avoid race condition. But it is wrong semantically.
//Data races on outLen also cause output[outLen++] to have data races.
//
//Data race pairs (we allow two pairs to preserve the original code pattern):
//1. outLen@34:9:W vs. outLen@34:9:W
//2. output[]@33:9:W vs. output[]@33:9:W

#include <omp.h>
#include <stdio.h>

int main()
{
    int i, inLen, outLen;
    int input[1000];
    int output[1000];

    inLen = 1000;
    outLen = 0;  // Adjust for 0-based indexing

    for (i = 0; i < inLen; i++) {
        input[i] = i + 1;
    }

    #pragma omp parallel for
    for (i = 0; i < inLen; i++) {
        output[outLen] = input[i];
        outLen = outLen + 1;
    }

    printf("output(500) = %d\n", output[499]);

    return 0;
} 