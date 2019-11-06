
//my cat a01227885


//all the includes
#include <stdio.h>
#include <sys/types.h>
#include <unistd.h>
#include <string.h>
#include <sys/stat.h>
#include <fcntl.h>


//magic
int main(int argc, char *argv[])
{
	
	//first we try to open
    int fp = open(argv[1], O_RDONLY);
    if (fp == -1)
    {
        printf("Can´t open. \n");
        return 1;
    }
	// here is to seek then close
    int size = lseek(fp, sizeof(char), SEEK_END);
    close(fp);
    fp = open(argv[1], O_RDONLY);
    if (fp == -1)
    {
        printf("Can´t open. \n");
        return 1;
    }
    char buf[size];
    read(fp, buf, size);
    close(fp);
    buf[size - 1] = '\0';

    write(1, buf, strlen(buf));

    return 0;
}
