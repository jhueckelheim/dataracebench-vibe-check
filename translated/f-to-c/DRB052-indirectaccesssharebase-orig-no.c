/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This example is to mimic a memory access pattern extracted from an LLNL proxy app.
//Two pointers have distance of 12.
//They are used as base addresses of two arrays, indexed through an index set.
//The index set has no two indices with distance of 12.
//So there is no loop carried dependence. No data race pairs.

#include <omp.h>
#include <stdio.h>
#include <stdlib.h>

int main()
{
    int* indexSet;
    double* xa1;
    double* xa2;
    double* base;
    int N = 180;
    int i, idx1, idx2;

    xa1 = (double*)malloc(2025 * sizeof(double));
    xa2 = (double*)malloc(2025 * sizeof(double));
    base = (double*)malloc(2025 * sizeof(double));
    indexSet = (int*)malloc(180 * sizeof(int));
    
    xa1 = base;
    xa2 = base;

    int temp_indexSet[180] = {521, 523, 525, 527, 529, 531, 547, 549,
        551, 553, 555, 557, 573, 575, 577, 579, 581, 583, 599,
        601, 603, 605, 607, 609, 625, 627, 629, 631, 633, 635,
        651, 653, 655, 657, 659, 661, 859, 861, 863, 865, 867,
        869, 885, 887, 889, 891, 893, 895, 911, 913, 915, 917,
        919, 921, 937, 939, 941, 943, 945, 947, 963, 965, 967,
        969, 971, 973, 989, 991, 993, 995, 997, 999, 1197, 1199,
        1201, 1203, 1205, 1207, 1223, 1225, 1227, 1229, 1231,
        1233, 1249, 1251, 1253, 1255, 1257, 1259, 1275, 1277,
        1279, 1281, 1283, 1285, 1301, 1303, 1305, 1307, 1309,
        1311, 1327, 1329, 1331, 1333, 1335, 1337, 1535, 1537,
        1539, 1541, 1543, 1545, 1561, 1563, 1565, 1567, 1569,
        1571, 1587, 1589, 1591, 1593, 1595, 1597, 1613, 1615,
        1617, 1619, 1621, 1623, 1639, 1641, 1643, 1645, 1647,
        1649, 1665, 1667, 1669, 1671, 1673, 1675, 1873, 1875,
        1877, 1879, 1881, 1883, 1899, 1901, 1903, 1905, 1907,
        1909, 1925, 1927, 1929, 1931, 1933, 1935, 1951, 1953,
        1955, 1957, 1959, 1961, 1977, 1979, 1981, 1983, 1985,
        1987, 2003, 2005, 2007, 2009, 2011, 2013};

    for (i = 0; i < 180; i++) {
        indexSet[i] = temp_indexSet[i];
    }

    for (i = 521; i <= 2025; i++) {
        base[i-1] = 0.0;  // Adjust for 0-based indexing
    }

    #pragma omp parallel for private(idx1, idx2)
    for (i = 0; i < N; i++) {
        idx1 = indexSet[i] - 1;  // Adjust for 0-based indexing
        idx2 = indexSet[i] + 12 - 1;  // Adjust for 0-based indexing
        base[idx1] = base[idx1] + 1.0;
        base[idx2] = base[idx2] + 3.0;
    }

    for (i = 521; i <= 2025; i++) {
        if (base[i-1] == 4.0) {
            printf("i = %d base = %f\n", i, base[i-1]);
        }
    }

    free(base);
    // Note: xa1 and xa2 point to base, so don't free them separately

    return 0;
} 