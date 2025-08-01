/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Example of writing to a file. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main()
{
    int i, len;
    int a[1000];
    FILE* file;

    len = 1000;

    for (i = 0; i < len; i++) {
        a[i] = i + 1;
    }

    if (access("mytempfile.txt", F_OK) == 0) {
        file = fopen("mytempfile.txt", "a");
    } else {
        file = fopen("mytempfile.txt", "w");
    }

    if (file != NULL) {
        #pragma omp parallel for
        for (i = 0; i < len; i++) {
            #pragma omp critical
            fprintf(file, "%d\n", a[i]);
        }
        
        fclose(file);
        remove("mytempfile.txt");  // Delete file
    }

    return 0;
} 