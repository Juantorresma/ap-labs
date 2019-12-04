#include <stdio.h>
#include <stdlib.h>
#include "omp.h"
#include "logger.h"

void main()
{
	#pragma omp parallel
	{
		int ID = omp_get_thread_num();
		infof(" hello(%d) ", ID);
		infof(" world(%d) \n", ID);
	}
}
