/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Copyright (c) 2017-20, Lawrence Livermore National Security, LLC
and DataRaceBench project contributors. See the DataRaceBench/COPYRIGHT file for details.

SPDX-License-Identifier: (BSD-3-Clause)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//This example is modified version of nestable_lock.1.c example, OpenMP 5.0 Application Programming Examples.
//A nested lock can be locked several times. It doesn't unlock until you have unset it as many times as the
//number of calls to omp_set_nest_lock.
//incr_b is called at line 54 and line 59. So, it needs a nest_lock for p%b@35:5. No data race.

#include <omp.h>
#include <stdio.h>

typedef struct {
    int a;
    int b;
    omp_nest_lock_t lck;
} pair;

void incr_a(pair* p, int a)
{
    p->a = p->a + 1;
}

void incr_b(pair* p, int b)
{
    omp_set_nest_lock(&p->lck);
    p->b = p->b + 1;
    omp_unset_nest_lock(&p->lck);
}

int main()
{
    int a, b;
    pair p;
    
    p.a = 0;
    p.b = 0;
    omp_init_nest_lock(&p.lck);

    #pragma omp parallel sections
    {
        #pragma omp section
        {
            omp_set_nest_lock(&p.lck);
            incr_b(&p, a);
            incr_a(&p, b);
            omp_unset_nest_lock(&p.lck);
        }

        #pragma omp section
        {
            incr_b(&p, b);
        }
    }

    omp_destroy_nest_lock(&p.lck);

    printf("%d\n", p.b);

    return 0;
} 