#include <stdio.h>
#include "logger.h"

int initLogger(char *logType) {
    printf("Initializing Logger on: %s\n", logType);
    return 0;
}

int infof(const char *format, ...) {
    return 0;
}

int warnf(const char *format, ...) {
    return 0;
}

int errorf(const char *format, ...) {
    return 0;
}
int panicf(const char *format, ...) {
    return 0;
}


//My logger a01227885

#include <stdio.h>
#include <stdarg.h>
#include <signal.h>
#include <stdlib.h>
#include <sys/types.h> 
#include <unistd.h> 
#include "logger.h"

//Types of actions

#define RESET		0
#define BRIGHT 		1
#define DIM		    2
#define UNDERLINE 	3
#define BLINK		4
#define REVERSE		7
#define HIDDEN		8

//colors im using to make it look pretty


#define BLACK 		0
#define RED		    1
#define GREEN		2
#define YELLOW		3
#define BLUE		4
#define MAGENTA		5
#define CYAN		6
#define	WHITE		7

//here i initialize 
int initLogger(char *logType)
{
    if (strcmp(logType, "syslog") == 0)
    {
        logDest = SYSLOG;
        printf("Initializing Logger on: %s\n", logType);
    }
    else if (strcmp(logType, "stdout") == 0)
    {
        logDest = STDOUT;
        printf("Initializing Logger on: %s\n", logType);
    }
    else
    {
        printf("Invalid log destination");
    }
    return 0;
}



//for the color
void textcolor(int attrib, int fore, int back)
{	char command[13];

	/* Command is the control command to the terminal */
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attrib, fore + 30, back + 40);
	printf("%s", command);
}


int infof(const char *format, ...) {

    int done;
    va_list arg;
	va_start (arg, format);
    textcolor(BRIGHT, CYAN, BLACK);
    printf("INFO: ");
    done = vfprintf (stdout, format, arg);
    va_end (arg);
    textcolor(RESET, RED, BLACK);	
    return done;

}

int warnf(const char *format, ...) {

    int done;
    va_list arg;
	va_start (arg, format);
    textcolor(BRIGHT, WHITE, BLACK);
    printf("WARNING: ");
    done = vfprintf (stdout, format, arg);
    va_end (arg);
    textcolor(RESET, YELLOW, BLACK);	
    return done;

}

int errorf(const char *format, ...) {
    
    int done;
    va_list arg;
	va_start (arg, format);
    textcolor(BRIGHT, RED, BLACK);
    printf("ERROR: ");
    done = vfprintf (stdout, format, arg);
    va_end (arg);
    textcolor(RESET, GREEN, BLACK);	
    return done;

}

int panicf(const char *format, ...) {
    
    int done;
    va_list arg;
	va_start (arg, format);
    textcolor(BRIGHT, RED, BLACK);
    printf("PANIC: ");
    done = vfprintf (stdout, format, arg);
    va_end (arg);
    textcolor(RESET, BLUE, BLACK);
    kill(getpid(), SIGQUIT);
    return done;

}

int printWithFormat(char *type, int color, const char *format, va_list arg)
{
    int done;
    textcolor(BRIGHT, color, HIDDEN);
    printf("%s", type);
    done = vfprintf(stdout, format, arg);
    textcolor(RESET, BLACK, HIDDEN);
    return done;
}