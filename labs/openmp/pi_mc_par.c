#include <stdio.h>
#include <omp.h>
#include "logger.h"
#include "random.h"

//
// The monte carlo pi program
//

static long num_trials = 1000000;

int main ()
{
	long i; long Ncirc = 0; double pi, x, y;
	double r = 1.0;
	seed(-r, r);
	#pragma omp parallel for private (x, y) reduction (+:Ncirc)
	for(i=0;i<num_trials; i++)
	{
		x = random(); y = random();
		if (( x*x + y*y) <= r*r)
			Ncirc++;
	}
	pi = 4.0 * ((double)Ncirc/(double)num_trials);
	infof("\n %d trials, pi is %f \n",num_trials, pi);
}
