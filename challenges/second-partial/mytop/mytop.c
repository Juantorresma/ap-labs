//My program A01227885 mytop Second partial

#include <stdio.h>
#include <unistd.h>
#include <dirent.h>
#include <ctype.h>
#include <string.h>
#include <stdlib.h>

//Here i define my struct which i will use in my program

typedef struct node {
	int pid;
	int ppid;
	char *name;
	char *state;
	int memory;
	int threads;
	int open_files;
	struct node * next;
} node_t;
//The nodes i will be using, start my head in null
node_t * head = NULL;

//my refresh funtion to refresh the head
void refresh(node_t * head){
	node_t * current = head;
	node_t * tmp;
    while (current != NULL) {
		tmp = current->next;
		free(current);
        current = tmp;
    }
}


//My put function to put the information in
void put(int pid, int ppid, char *name, char *state, int memory, int threads, int open_files) {
	node_t * current = head;
	node_t * parent = NULL;
	int new_head = 1;
    while (current != NULL && current->memory > memory) {
		new_head = 0;
		parent = current;
        current = current->next;
    }
	node_t *tmp = current;
    current = malloc(sizeof(node_t));
	current->pid = pid;
	current->ppid = ppid;
	current->name = name;
	current->state = state;
	current->memory = memory;
	current->threads = threads;
	current->open_files = open_files;
    current->next = tmp;
	if(new_head){
		head = current;
	}else{
		parent->next = current;
	}
}

//My proc function here is the magic
int proc(){
	head = NULL;
	DIR *d;
	struct dirent *dir;
	d = opendir("/proc/");
	if (d)
	{
		while ((dir = readdir(d)) != NULL)
		{
			if (isdigit((dir->d_name)[0])) {
				char *path = (char *)malloc(sizeof(char)*(6+strlen(dir->d_name)+8));
				strcat(path, "/proc/");
				strcat(path, dir->d_name);
				strcat(path, "/status");
				FILE *fp = fopen(path, "r");
				char * line = NULL;
			    size_t len = 0;
			    ssize_t read;
				if (fp == NULL){
					printf("Can't open file %s\n", path);
					return -1;
				}
				int pid;
				int ppid;
				char *name;
				char *state;
				int memory;
				int threads;
				int open_files;
				char *substr;
				char *p;
				while ((read = getline(&line, &len, fp)) != -1) {
					if (line[read - 1] == '\n')
				    {
				        line[read - 1] = '\0';
				        --read;
				    }
					// Pid
					substr = "Pid:\t";
					p = strstr(line, substr);
					if(p == line) {
						pid = atoi((p+5));
					}
					// Parent id
					substr = "PPid:\t";
					p = strstr(line, substr);
					if(p == line) {
						ppid = atoi((p+6));
					}
					// Name
					substr = "Name:\t";
					p = strstr(line, substr);
					if(p == line) {
						p = (p+6);
						name = (char *)malloc(sizeof(char)*(strlen(p)+1));
						strcpy(name, p);
					}
					// State
					substr = "State:\t";
					p = strstr(line, substr);
					if(p == line) {
						p = (p+10);
						p[strlen(p) - 1] = '\0';
						state = (char *)malloc(sizeof(char)*(strlen(p)));
						strcpy(state, p);
					}
					// Memory
					substr = "VmSize:\t";
					p = strstr(line, substr);
					if(p == line) {
						p = (p+8);
						p[strlen(p) - 3] = '\0';
						memory = atoi(p)/1024;
					}
					// The number of hreads
					substr = " # Threads:\t";
					p = strstr(line, substr);
					if(p == line) {
						threads = atoi((p+9));
					}

				}
				fclose(fp);
				free(path);

				char *fdpath = (char *)malloc(sizeof(char)*(6+strlen(dir->d_name)+5));
				strcat(fdpath, "/proc/");
				strcat(fdpath, dir->d_name);
				strcat(fdpath, "/fd/");
				DIR *procfd;
				struct dirent *fd;
				procfd = opendir(fdpath);
				if (procfd)
				{
					open_files = -2;
					while ((fd = readdir(procfd)) != NULL)
					{
						open_files++;
					}
				}
				closedir(procfd);

				put(pid, ppid, name, state, memory, threads, open_files);

				free(fdpath);
				if (line)
		 		   free(line);
			}
		}
		closedir(d);
	}
	return 0;
}


//clear the screen, i had to modify it
void clear() {
	for(int i = 0; i < 22; i++)
    	printf("\033[A\r");
}

//the ptinting i divided so the header would only print once
void print_header(){
	printf("| %7s ", "PID");
	printf("| %7s ", "Parent");
	printf("| %36s ", "Name");
	printf("| %8s ", "State");
	printf("| %8s ", "Memory");
	printf("| %9s ", "# Threads");
	printf("| %10s |\n", "Open Files");
	printf("|---------|---------|--------------------------------------|----------|----------|-----------|------------|\n");
}

//then every node
void print_node(node_t *current){
	printf("| %7d ", current->pid);
	printf("| %7d ", current->ppid);
	printf("| %36s ", current->name);
	printf("| %8s ", current->state);
	printf("| %7dM ", current->memory);
	printf("| %9d ", current->threads);
	printf("| %10d |\n", current->open_files);
}

//here i put together my printt funcions
void print(){
	print_header();
  	node_t * current = head;
	int i = 0;
	while (current != NULL && i<20) {
		i++;
	    print_node(current);
	    current = current->next;
	}
	refresh(head);
	sleep(2);
	clear();
}

//infinite loop to call my process
int main(){
	while(1){
		int out = proc();
		if(out!=0){
			return -1;
		}
		print();
	}
	return 0;
}
