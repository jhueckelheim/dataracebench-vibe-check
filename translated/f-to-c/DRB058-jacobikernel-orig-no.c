/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//Two parallel for loops within one single parallel region,
//combined with private() and reduction().

//3.7969326424804763E-007 vs 3.7969326424804758E-007. There is no race condition. The minute
//difference at 22nd point after decimal is due to the precision in fortran95

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>
#include <math.h>

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
    tol = 0.0000000001;
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

    for (i = 0; i < n; i++) {
        for (j = 0; j < m; j++) {
            xx = (int)(-1.0 + dx * i);
            yy = (int)(-1.0 + dy * i);
            u[i][j] = 0.0;
            f[i][j] = -1.0 * alpha * (1.0 - xx*xx) * (1.0 - yy*yy) - 2.0 * (1.0 - xx*xx) - 2.0 * (1.0 - yy*yy);
        }
    }
}

void jacobi()
{
    double omega;
    int i, j, k;
    double error, resid, ax, ay, b;

    MSIZE = 200;
    mits = 1000;
    tol = 0.0000000001;
    relax = 1.0;
    alpha = 0.0543;
    n = MSIZE;
    m = MSIZE;

    omega = relax;
    dx = 2.0 / (double)(n - 1);
    dy = 2.0 / (double)(m - 1);

    ax = 1.0 / (dx * dx);         // X-direction coef
    ay = 1.0 / (dy * dy);         // Y-direction coef
    b = -2.0 / (dx * dx) - 2.0 / (dy * dy) - alpha;

    error = 10.0 * tol;
    k = 1;

    for (k = 0; k < mits; k++) {
        error = 0.0;

        // Copy new solution into old
        #pragma omp parallel
        {
            #pragma omp for private(i, j)
            for (i = 0; i < n; i++) {
                for (j = 0; j < m; j++) {
                    uold[i][j] = u[i][j];
                }
            }
            
            #pragma omp for private(i, j, resid) reduction(+:error)
            for (i = 1; i < (n-1); i++) {
                for (j = 1; j < (m-1); j++) {
                    resid = (ax * (uold[i-1][j] + uold[i+1][j]) + ay * (uold[i][j-1] + uold[i][j+1]) + b * uold[i][j] - f[i][j]) / b;
                    u[i][j] = uold[i][j] - omega * resid;
                    error = error + resid * resid;
                }
            }
        }

        // Error check
        error = sqrt(error) / (n * m);
    }

    printf("Total number of iterations: %d\n", k);
    printf("Residual: %f\n", error);
}

int main()
{
    initialize();
    jacobi();

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