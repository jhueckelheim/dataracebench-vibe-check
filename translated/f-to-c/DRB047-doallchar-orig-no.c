/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//One dimension array computation
//with finer granularity than traditional 4 bytes.
//There is no data race pair.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main()
{
    char (*a)[101];  // Array of strings with length 100+1 for null terminator
    char str[51];
    int i;

    a = (char (*)[101])malloc(100 * 101 * sizeof(char));

    #pragma omp parallel for private(str)
    for (i = 0; i < 100; i++) {
        sprintf(str, "%10d", i + 1);  // Adjust for 1-based indexing in output
        strcpy(a[i], str);
    }

    printf("a(i) %s\n", a[22]);  // Adjust for 0-based indexing

    free(a);

    return 0;
} 