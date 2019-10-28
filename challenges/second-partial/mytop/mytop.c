// Fernando Martinez
// A01630401

#include <stdio.h>
#include <dirent.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/inotify.h>

#define EVENT_SIZE (sizeof(struct inotify_event))
#define EVENT_BUF_LEN (1024*(EVENT_SIZE + 16))



void clear() {
	printf("\e[1;1H\e[2J");
}


void checkChanges(){
	int fd;
        int wd;
        int rd;
        int p = 0;
        struct inotify_event *event;
        char buffer[EVENT_BUF_LEN];
        fd = inotify_init();
        if(fd<0)
                printf("inotify error\n");
        wd = inotify_add_watch(fd, "/proc", IN_CREATE | IN_DELETE);
        while(1) {
                rd = read(fd, buffer, EVENT_BUF_LEN);
                if(read < 0)
                        printf("event not read\n");
                while(p < rd) {
                        event = (struct inotify_event * ) &buffer[p];
                        if(event->len) {
                                clear();
                                getProcesses();
                                p += EVENT_SIZE + event->len;
                        }
                }
        }
        inotify_rm_watch(fd, wd);
        close(fd);
}

void getProcesses() {
	char procContent[200][15];
	int lenAll = 0;
	DIR *directory;
	struct dirent *entry;
	if((directory = opendir("/proc")) != NULL) {
		while((entry = readdir(directory)) != NULL) {
			strcpy(procContent[lenAll], entry->d_name);
			lenAll++;
		}
		closedir(directory);
	}
	else
		printf("Directory not open");


	char allProcesses[1000][30];
	int i;
	int q;
	int len2 = 0;
	char data[10000];
	for(i = 0; i <= lenAll; i++) {
		if(atoi(procContent[i]) != 0) {
			strcpy(allProcesses[len2], "/proc/");
			strcat(allProcesses[len2], procContent[i]);
			//strcat(allProcesses[len2], "/status");
			//printf("%s\n", allProcesses[len2]);
			len2++;
		}
	}
	printf("| PID \t| Parent| Name                        | State    | Memory    | # Threads | Open Files |\n");
	printf("|-------|-------|-----------------------------|----------|-----------|-----------|------------|\n");
	for(q = 0; q < len2; q++) {
		char fileName[30];
		char files[30];
		strcpy(fileName, allProcesses[q]);
		strcat(fileName, "/status");
		strcpy(files, allProcesses[q]);
		strcat(files, "/fd/");
		int fileDesc = open(fileName, O_RDONLY);
		//printf("Name of the file:                  %s\n", allProcesses[q]);
		if(fileDesc < 0) {
			printf("File not found\n");
			//printf("%s\n", allProcesses[i]);
		}
		else {
			//printf("File read\n");
			//printf("%s\n", allProcesses[i]);
			//printf("\n");
			int fileRead = read(fileDesc, data, 10000);
			if(fileRead < 0) {
				printf("File not read");
			}
			close(fileDesc);
			data[fileRead] = '\0';
			int u;
			char pid[15];
			char parent[15];
			char name[30];
			char status[20];
			char memory[30];
			char threads[15];
			int oFiles = 0;
			char *line;
			char *toSearch;

			toSearch = "Pid:\t";
			line = strstr(data, toSearch);
			if(line != NULL) {
				u = 0;
				line += 5;
				while(line[u] != '\n') {
					pid[u] = line[u];
					u++;
				}
			}

			toSearch = "PPid:\t";
			line = strstr(data, toSearch);
			if(line != NULL) {
				u = 0;
				line += 6;
				while(line[u] != '\n') {
					parent[u] = line[u];
					u++;
				}
			}

			toSearch = "Name:\t";
			line = strstr(data, toSearch);
			if(line != NULL) {
				u = 0;
				line += 6;
				while(line[u] != '\n'){
					name[u] = line[u];
					u++;
				}
			}

			toSearch = "State:\t";
			line = strstr(data, toSearch);
			if(line != NULL) {
				u = 0;
				line += 10;
				while(line[u] != ')') {
					status[u] = line[u];
					u++;
				}
			}


			toSearch = "VmSize:\t";
			line = strstr(data, toSearch);
			if(line != NULL) {
				u = 0;
				line += 10;
				while(line[u] != '\n') {
					memory[u] = line[u];
					u++;
				}
			}

			toSearch = "Threads:\t";
			line = strstr(data, toSearch);
			if(line != NULL) {
				u = 0;
				line += 9;
				while(line[u] != '\n') {
					threads[u] = line[u];
					u++;
				}
			}

			DIR *dirFiles;
			if((dirFiles = opendir(files))!= NULL) {
				while((readdir(dirFiles)) != NULL)
					oFiles++;
				closedir(dirFiles);
			}


			printf("| %s \t| %s \t| %27s | %8s | %9s | %9s | %10d |\n", pid, parent, name, status, memory, threads, oFiles);
			int m;
			for(m = 0; m < 50; m++) {
				name[m] = '\0';
				status[m] = '\0';
				pid[m] = '\0';
				parent[m] = '\0';
				memory[m] = '\0';
				threads[m] = '\0';
                	}
		}

	}
}

int main(){
	clear();
	getProcesses();
	checkChanges();
	return 0;
}

