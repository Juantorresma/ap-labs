#include <stdio.h>

/* filecopy:  copy file ifp to file ofp */
void filecopy(FILE *ifp, FILE *ofp)
{
    int c;

    while ((c = getc(ifp)) != EOF)
        putc(c, ofp);

}

/* cat:  concatenate files, version 2 */
int main(int argc, char *argv[])
{
    FILE *fp;
    void filecopy(FILE *, FILE *);
    char *prog = argv[0];   /* program name for errors */

    if (argc == 1)  /* no args; copy standard input */
        filecopy(stdin, stdout);
    else
        while (--argc > 0)
            if ((fp = fopen(*++argv, "r")) == NULL) {
                fprintf(stderr, "%s: can′t open %s\n",
			prog, *argv);
                return 1;
            } else {
                filecopy(fp, stdout);
                fclose(fp);
            }

    if (ferror(stdout)) {
        fprintf(stderr, "%s: error writing stdout\n", prog);
        return 2;
    }

    return 0;
}

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
    close(fd);
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
