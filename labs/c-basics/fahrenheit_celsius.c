//A01227885

#include <stdio.h>

/* print Fahrenheit-Celsius table */

int main(int argc, char **argv)
{
    int fahr;
    int start;
    int end;
    int increment;
    
    if (argc == 2){
        fahr =  atoi(argv[1]);
        printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
    }
    
    
     if (argc == 4){
         start =  atoi(argv[1]);
         end =  atoi(argv[2]);
         increment =  atoi(argv[3]) ;
         
         
        for (fahr = start; fahr <= end; fahr = fahr + increment)
	    printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
     }
    return 0;
}


