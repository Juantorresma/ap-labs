#include <stdio.h>
#include <omp.h>
#include "logger.h"

static long num_steps = 1000000;
#define NUM_THREADS 2
double step;
int main ()
{
    double x, pi, sum = 0.0;
    double start_time, run_time;

    step = 1.0/(double) num_steps;

    start_time = omp_get_wtime();
    omp_set_num_threads(NUM_THREADS);
	#pragma omp parallel for private(x)reduction(+:sum)

    for (int i=1;i<= num_steps; i++){
	x = (i-0.5)*step;
	sum = sum + 4.0/(1.0+x*x);
    }

    pi = step * sum;
    run_time = omp_get_wtime() - start_time;
    infof("\n pi with %d steps is %f in %f seconds ",num_steps,pi,run_time);
}
