//My program A01227885 mytop Second partial

#include <stdio.h>
#include <unistd.h>
#include <dirent.h>
#include <ctype.h>
#include <string.h>
#include <stdlib.h>

//Here i define my dtruct which i will use in my program

typedef struct node {
	int process_id;
	int parent_id;
	char *name;
	char *state;
	int memory;
	int num_threads;
	int open_files;
	struct node * next;
} 

//The nodes i will be using, start my headd in null
node_t;
node_t * head = NULL;

//My push function
void push(int process_id, int parent_id, char *name, char *state, int memory, int num_threads, int open_files) {
	node_t * current = head;
	node_t * parent = NULL;
	int newHead = 1;
    while (current != NULL && current->memory > memory) {
		newHead = 0;
		parent = current;
        current = current->next;
    }
	//Assign variables
	
	node_t *tmp = current;
	current = malloc(sizeof(node_t));
	current->process_id = process_id;
	current->parent_id = parent_id;
	current->name = name;
	current->state = state;
	current->memory = memory;
	current->num_threads = num_threads;
	current->open_files = open_files;
	current->next = tmp;
	if(newHead){
		head = current;
	}else{
		parent->next = current;
	}
}

//My clean funtion
void clean(node_t * head){
	node_t * current = head;
	node_t * tmp;
    while (current != NULL) {
		tmp = current->next;
		free(current);
        current = tmp;
    }
}

//My proc function
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
				int process_id;
				int parent_id;
				char *name;
				char *state;
				int memory;
				int num_threads;
				int open_files;
				char *substr;
				char *p;
				while ((read = getline(&line, &len, fp)) != -1) {
					if (line[read - 1] == '\n')
				    {
				        line[read - 1] = '\0';
				        --read;
				    }
					// Process id
					substr = "PID:\t";
					p = strstr(line, substr);
					if(p == line) {
						process_id = atoi((p+5));
					}
					// Parent id
					substr = "Parent:\t";
					p = strstr(line, substr);
					if(p == line) {
						parent_id = atoi((p+6));
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
					// # Threads
					substr = "# Threads:\t";
					p = strstr(line, substr);
					if(p == line) {
						num_threads = atoi((p+9));
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
				push(process_id, parent_id, name, state, memory, num_threads, open_files);
				free(fdpath);
				if (line)
		 		   free(line);
			}
		}
		closedir(d);
	}
	return 0;
}


//The clear function
void clear() {
	for(int i = 0; i < 22; i++)
    	printf("\033[A\r");

}

//Here i print the headers
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

//Here i ptint the current node
void print_node(node_t *current){
	printf("| %7d ", current->process_id);
	printf("| %7d ", current->parent_id);
	printf("| %36s ", current->name);
	printf("| %8s ", current->state);
	printf("| %7dM ", current->memory);
	printf("| %9d ", current->num_threads);
	printf("| %10d |\n", current->open_files);
}


//My print function which will call print node and header
void print(){
	print_header();
  	node_t * current = head;
	int i = 0;
	while (current != NULL && i<20) {
		i++;
	    print_node(current);
	    current = current->next;
	}
	clean(head);
	sleep(2);
	clear();
}

//An infinite loop to run my process proc
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
