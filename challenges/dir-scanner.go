//Primer Parcial A01227885

//primero incluimos todas las librerias que vamos a necesitar
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <stdbool.h>


//Definimos los status de nuestros paquetes
#define installed 0
#define removed 1
#define upgraded 2


//Asi mismo definimos las acciones que haremos
#define fetchDate 10
#define findAction 11
#define fetchAction 12
#define fetchName 13
#define findLine 14


//en este archivo se guardara el reporte
#define reportFile "packages_report.txt"

//Aqui generamos un paquete para transportar y utilizar estos datos juntos
struct Package
{
    char packName[50];
    char packDate[17];
    char packUpdate[17];
    int numUpdates;
    char removalDate[17];
    int currentStatus;
};

//Aqui genero mi hashtable, con tamaño elementos y un arreglo de mis paquetes
struct Hashtable
{
    int size;
    int numElements;
    struct Package array[1000];
};

//aqui convertimos los paquetes a strings con el formato dado

void pToString(char string[], struct Package *pack)
{
    strcat(string, "- Package Name        : ");
    strcat(string, pack->packName);
    strcat(string, "\n");
    strcat(string, "  - Install date      : ");
    strcat(string, pack->packDate);
    strcat(string, "\n");
    strcat(string, "  - Last update date  : ");
    strcat(string, pack->packUpdate);
    strcat(string, "\n");
    strcat(string, "  - How many updates  : ");
    char numBuf[20];
    sprintf(numBuf, "%d\n", pack->numUpdates);
    strcat(string, numBuf);
    strcat(string, "  - Removal date      : ");
    strcat(string, pack->removalDate);
}

//aqui convertimos las Hashtables a strings con el formato dado

void htToString(char string[], struct Hashtable *table)
{
    for (int i = 0; i < table->size; i++)
    {
        if (strcmp(table->array[i].packName, "") != 0)
        {
            pToString(string, &table->array[i]);
            strcat(string, "\n\n");
        }
    }
}


//esta funcion probara si es una accion lo que este en nuestro buffer
bool isAction(char char1, char char2, char char3)
{
    if (char1 == 'i' && char2 == 'n' && char3 == 's')
    {
        return true;
    }
    else if (char1 == 'u')
    {
        return true;
    }
    else if (char1 == 'r' && char2 == 'e')
    {
        return true;
    }
    else
    {
        return false;
    }
}

//Esta es la funcion para sacar el codigo hash
int getHashCode(char code[])
{
    int len = strlen(code);
    int hashValue = 0;

    for (int i = 0; i < len; i++)
    {
        hashValue = hashValue * 31 + code[i];
    }

    hashValue = hashValue & 0x7fffffff;
    return hashValue;
}


//Esta solo añade el paquete a la hashtable
void addToHashtable(struct Hashtable *table, struct Package *pack)
{
    for (int i = 0; i < table->numElements + 1; i++)
    {
        int hashValue = getHashCode(pack->packName) + i;
        int index = hashValue % table->size;
        if (strcmp(table->array[index].packName, "") == 0)
        {
            table->array[index] = *pack;
            break;
        }
    }

    //se incrementa el numero de elementos en la hash
    table->numElements += 1;
}


// para encontrar una linea nueva en el hash table
bool findInHashtable(struct Hashtable *table, char key[])
{
    for (int i = 0; i < table->numElements + 1; i++)
    {
        int hashValue = getHashCode(key) + i;
        int index = hashValue % table->size;
        if (strcmp(table->array[index].packName, key) == 0)
        {
            return true;
        }
        else if (strcmp(table->array[index].packName, "") == 0)
        {
            return false;
        }
    }
    return false;
}

//Esta funcion nos servira para obtener los paquetes
struct Package *get(struct Hashtable *table, char key[])
{
    for (int i = 0; i < table->numElements + 1; i++)
    {
        int hashValue = getHashCode(key) + i;
        int index = hashValue % table->size;
        if (strcmp(table->array[index].packName, key) == 0)
        {
            return &table->array[index];
        }
        else if (strcmp(table->array[index].packName, "") == 0)
        {
            return NULL;
        }
    }
    return NULL;
}


//para generar el reporte con el formato requerido  primero ponemos el titulo y despues hacemos el htToString para imprimirtoda la informacin de nuestro reporte que estara en la hash table de paquetes
void makeReport(char *reportS, int installedPack, int removedPack, int upgradedPack, int currentPack, struct Hashtable *table)
{
    strcat(reportS, "Pacman Packages Report\n");
    strcat(reportS, "----------------------\n");
    char numBuf[120];
    sprintf(numBuf, "- Installed packages : %d\n- Removed packages   : %d\n- Upgraded packages  : %d\n- Current installed  : %d\n\n", installedPack, removedPack, upgradedPack, currentPack);
    strcat(reportS, numBuf);
    htToString(reportS, table);
}

//Aqui pasa la mágia, necesitaremos nuetsro archivo con los logs y el reporte
void analizeLog(char *logFile, char *report)
{
    //avisamos de que archivo (log) se generara el paquete
    printf("The report is being generated from the log file [%s] \n", logFile);

    //inicializamos nuestra hash table vacia con tamaño de mil, al igual que nuestras variables en 0 ya que seran contadores

    struct Hashtable table = {1000, 0};
    int installedPack = 0;
    int removedPack = 0;
    int upgradedPack = 0;
    int currentPack = 0;

    //intentamos abrir el logfile 
    int fd = open(logFile, O_RDONLY);
    if (fd == -1)
    {
        printf("The file could not be opened.\n");
        return;
    }
    //definimos el tamaño
    int size = lseek(fd, sizeof(char), SEEK_END);
    close(fd);
    fd = open(logFile, O_RDONLY);
    if (fd == -1)
    {
        printf("The file could not be opened.\n");
        return;
    }
    //inicializamos nuestro buffer de caracteres tamaño size
    char buf[size];
    read(fd, buf, size);
    close(fd);
    //llenamos la ultima posicion de nuestro buffer
    buf[size - 1] = '\0';

    //inicializamos variables y el estado en la fecha de obtencion, para ir cambiando este ultimo con un case
    int i = 0;
    int j = 0;
    int state = fetchDate;
    char date[17];
    char name[50];
    char action[10];
    bool isValid = false;

    while (i < size)
    {
        switch (state)
        {
        //definimos la fecha y al hacerlo pasamos a encontrar la siguiente accion
        case fetchDate:
            if (buf[i] != 'f')
            {
                i++;
                j = 0;
                while (buf[i] != ']')
                {
                    date[j] = buf[i];
                    j++;
                    i++;
                }
                date[j] = '\0';
                i = i + 2;
                state = findAction;
            }
            else
            {
                state = findAction;
            }
            break;
        //al encontrar la siguiente accion la definimos
        case findAction:
            while (buf[i] != ' ')
            {
                i++;
            }
            i++;
            state = fetchAction;
            break;
        //definimos nuestra nuesva accion mandando llamar la funcion para saber si es una accion, en caso de encontrarla definimos el nombre, sino buscamos la siguiente línea
        case fetchAction:
            j = 0;
            if (isAction(buf[i], buf[i + 1], buf[i+2]))
            {
                isValid = true;
                while (buf[i] != ' ')
                {
                    action[j] = buf[i];
                    i++;
                    j++;
                }
                action[j] = '\0';
                i++;
                state = fetchName;
            }
            else
            {
                state = findLine;
            }
            break;
        //aqui definimos el nombre y buscamos nueva linea
        case fetchName:
            j = 0;
            while (buf[i] != ' ')
            {
                name[j] = buf[i];
                i++;
                j++;
            }
            name[j] = '\0';
            i++;
            state = findLine;
            break;
        //tenemos que buscar si existe en la tabla hash, si no lo ponemos 
        case findLine:
            while (!(buf[i] == '\n' || buf[i] == '\0'))
            {
                i++;
            }
            i++;
            if (isValid)
            {
                if (!findInHashtable(&table, name))
                {
                    struct Package pack = {"", "", "", 0, "-", installed};
                    strcpy(pack.packName, name);
                    strcpy(pack.packDate, date);
                    addToHashtable(&table, &pack);

                    installedPack++;
                }
                else
                {
                    struct Package *p1 = get(&table, name);
                    if (strcmp(action, "installed") == 0)
                    {
                        if (p1->currentStatus == removed)
                        {
                            p1->currentStatus = installed;
                            strcpy(p1->removalDate, "-");
                            removedPack--;
                        }
                    }
                    else if (strcmp(action, "removed") == 0)
                    {
                        if (p1->currentStatus == installed || p1->currentStatus == upgraded)
                        {
                            p1->currentStatus = removed;
                            strcpy(p1->removalDate, date);
                            strcpy(p1->packUpdate, date);
                            p1->numUpdates = p1->numUpdates + 1;
                            removedPack++;
                        }
                    }
                    else if (strcmp(action, "upgraded") == 0)
                    {
                        if (p1->currentStatus == installed)
                        {
                            p1->currentStatus = upgraded;
                            strcpy(p1->packUpdate, date);
                            p1->numUpdates = p1->numUpdates + 1;
                            upgradedPack++;
                        }
                        else if (p1->currentStatus == upgraded)
                        {
                            strcpy(p1->packUpdate, date);
                            p1->numUpdates = p1->numUpdates + 1;
                        }
                    }
                }
            }
            isValid = false;
            //regresamos al status inicial
            state = fetchDate;
            if (i >= size - 1)
            {
                i = i + 1;
            }
            break;
        }
    }
    //realizamos calculos y creamos nuestro reporte
    currentPack = installedPack - removedPack;
    char reportS[100000];
    makeReport(reportS, installedPack, removedPack, upgradedPack, currentPack, &table);
    
    fd = open(report, O_CREAT | O_WRONLY, 0600);
    if (fd == -1)
    {
        printf("The file could not be opened.\n");
        return;
    }
    write(fd, reportS, strlen(reportS));
    close(fd);
    //
    printf("Your report: [%s]\n", report);
}


//Aqui solo verificamos el número de argumentos y regresamos error sino es válido
int main(int argc, char **argv)
{

    if (argc < 2)
    {
        printf("Usage:./pacman-analizer.o pacman.txt\n");
        return 1;
    }

    analizeLog(argv[1], reportFile);

    return 0;
}

