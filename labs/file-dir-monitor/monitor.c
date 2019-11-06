#include <stdio.h>
#include "logger.h"
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/inotify.h>

#define BUF_LEN sizeof(struct inotify_event) * 1024

int var1;
int var2;
int var3;
char* point;
struct inotify_event *event;

void monitor(struct inotify_event* event){
     if (event->mask & IN_ACCESS)    infof("Acceso a archivo o a directorio.\n");
     if (event->mask & IN_CREATE)        warnf("Creación de archivo o directorio.\n");
     if (event->mask & IN_DELETE)        warnf("Eliminación de archivo o directorio.\n");
     if (event->mask & IN_OPEN)          infof("Se ha abierto un archivo o directorio.\n");
     if (event->mask & IN_MODIFY)          warnf("Se ha modificado un archivo o directorio.\n");
}

int main(int argc, char** argv){
    if(argc < 2){
        printf("Porfavor especifique un directorio para monitorear\n");
        return -1;
    }
    var1 = inotify_init1(O_NONBLOCK);
    var2 = inotify_add_watch(var1, argv[1], IN_ALL_EVENTS);
    char* buff = (char*)malloc(BUF_LEN);
    while(1){
        var3 = read(var1, buff, BUF_LEN);
        point = buff;
        event = (struct inotify_event*)point;
        for (point = buff; point < buff + var3; ) {
             event = (struct inotify_event *) point;
             monitor(event);
             point += sizeof(struct inotify_event) + event->len;
         }
    }
    close(var1);
    return 0;
}


