/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Use of private() clause. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

// Global variables (from module)
int MSIZE;
int n, m, mits;
double** u;
double** f;
double** uold;
double dx, dy, tol, relax, alpha;

void initialize()
{
    int i, j, xx, yy;

    MSIZE = 200;
    mits = 1000;
    relax = 1.0;
    alpha = 0.0543;
    n = MSIZE;
    m = MSIZE;

    // Allocate 2D arrays
    u = (double**)malloc(MSIZE * sizeof(double*));
    f = (double**)malloc(MSIZE * sizeof(double*));
    uold = (double**)malloc(MSIZE * sizeof(double*));
    for (i = 0; i < MSIZE; i++) {
        u[i] = (double*)malloc(MSIZE * sizeof(double));
        f[i] = (double*)malloc(MSIZE * sizeof(double));
        uold[i] = (double*)malloc(MSIZE * sizeof(double));
    }

    dx = 2.0 / (double)(n - 1);
    dy = 2.0 / (double)(m - 1);

    // Initialize initial condition and RHS
    #pragma omp parallel for private(i, j, xx, yy)
    for (i = 0; i < n; i++) {
        for (j = 0; j < m; j++) {
            xx = (int)(-1.0 + dx * i);
            yy = (int)(-1.0 + dy * i);
            u[i][j] = 0.0;
            f[i][j] = -1.0 * alpha * (1.0 - xx*xx) * (1.0 - yy*yy) - 2.0 * (1.0 - xx*xx) - 2.0 * (1.0 - yy*yy);
        }
    }
}

int main()
{
    initialize();

    // Free allocated memory
    for (int i = 0; i < MSIZE; i++) {
        free(u[i]);
        free(f[i]);
        free(uold[i]);
    }
    free(u);
    free(f);
    free(uold);

    return 0;
} 