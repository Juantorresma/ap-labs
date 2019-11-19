//Multiplier function A01227885
//Includes
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <strings.h>
#include "logger.h"

//define matrix sizes and number of elements
#define MATRIX_DIMENSION 2000
#define MATRIX_NUM_ELEMENTS 4000000

//functions that will be needed (Teacher)
long * readMatrix(char *filename);
long * getColumn(int col, long *matrix);
long * getRow(int row, long *matrix);
int getLock();
int releaseLock(int lock);
long dotProduct(long *vec1, long *vec2);
long * multiply(long *matA, long *matB);
int saveResultMatrix(long *result);

//Extra function needed for my code
void *threadFunc(void *arg);

//Globar variables (Teacher)
int NUM_BUFFERS;
long **buffers;
pthread_mutex_t *mutexes;
long *result;

//Global variable
pthread_t threads[MATRIX_DIMENSION];

//I make a struct to put data together
struct vectorStruct {
	int positionR;
	int positionC;
	long *multResult;
    long *matA;
	long *matB;
};


//functions
long * readMatrix(char *filename){
    
    FILE *f;

    if ((f = fopen(filename,"r")) == NULL){
        errorf("Error: could not open file");
        exit(1);
    }

    long *res = (long *)malloc(MATRIX_NUM_ELEMENTS * sizeof(long));
    char *num = NULL;
	int i;
	size_t len = 0;

    for(i = 0; (getline(&num, &len, f)) != -1 ; i++)
        res[i] = strtol(num, NULL, 10);

    fclose(f); 
	free(num);
    return res;
}

long * getColumn(int col, long *matrix){

    if (col < 0 || col > MATRIX_DIMENSION) {
		warnf("The colum has not a valid number, it must be positive and below 2000\n");
		exit(1);
	}

    long pos = col - 1;
	long *vector;
	vector = (long *)malloc(MATRIX_DIMENSION * sizeof(long));

	for (int i = 0; i < MATRIX_DIMENSION; i++)
		vector[i] = matrix[pos];
        pos += MATRIX_DIMENSION;
	
	return vector;
}

long * getRow(int row, long *matrix){

    if (row < 0 || row > MATRIX_DIMENSION) {
		warnf("The row has not a valid number, it must be positive and below 2000\n");
		exit(1);
	}

    int pos = ((2 * row) - 2) * 1000;
    long *vector;
	vector = (long *)malloc(MATRIX_DIMENSION * sizeof(long));
    
    for (int i = 0; i < MATRIX_DIMENSION; i++) 
	    vector[i] = matrix[pos++];
	
	return vector;
}

int getLock(){
    for (int i = 0; i < NUM_BUFFERS; i++) {
		if (pthread_mutex_trylock(&mutexes[i]) == 0)
			return i;
	}
	return -1;
}

int releaseLock(int lock){
    return ((pthread_mutex_unlock(&mutexes[lock]) == 0)? 0 : -1 );
}

long dotProduct(long *vec1, long *vec2){
    long res = 0;
	int i;
    for ( i = 0; i < 2000; i++)
		res += vec1[i] * vec2[i];
	
    return res;
}

long * multiply(long *matA, long *matB){

    long *result = (long *)malloc(MATRIX_NUM_ELEMENTS * sizeof(long));
	
    int i, 
        j;

	for ( i = 0; i < MATRIX_DIMENSION; i++) {
		for ( j = 0; j < MATRIX_DIMENSION; j++) {
			struct vectorStruct *vector;
            vector =(struct vectorStruct *) malloc(sizeof(struct vectorStruct));

			vector->matA = matA;
			vector->matB = matB;
			vector->positionR = i + 1;
			vector->positionC = j + 1;
			vector->multResult = result;

			pthread_create(&threads[j], NULL , threadFunc , (void *)vector);
		}

		for ( j = 0; j < MATRIX_DIMENSION; j++)
			pthread_join(threads[j], NULL);

		fflush(stdout);
	}

	return result;
}

int saveResultMatrix(long *result){

    FILE *f = fopen("result.dat", "w");
	if (f == NULL) {
		errorf("Could not write the file 'result.dat' \n");
		return -1;
	}

    long i;
	for ( i = 0; i < MATRIX_NUM_ELEMENTS; i++)
		fprintf(f, "%ld\n", result[i]);

	fclose(f);
	return 0;
}

void *threadFunc(void *arg)
{
	struct vectorStruct *currVec = (struct vectorStruct *)arg;
    long i;
	int lockRow, lockCol;

	while ((lockRow = getLock()) == -1);
	while ((lockCol = getLock()) == -1);

	buffers[lockRow] = getRow(currVec->positionR, currVec->matA);
	buffers[lockCol] = getColumn(currVec->positionC, currVec->matB);
	
	i = ((((currVec->positionR - 1) * MATRIX_DIMENSION) + currVec->positionC) - 1);
	currVec->multResult[i] = dotProduct(buffers[lockRow], buffers[lockCol]);

	free(buffers[lockRow]);
	free(buffers[lockCol]);
	free(arg);

	while (releaseLock(lockRow) != 0);
	while (releaseLock(lockCol) != 0);

	return NULL;
}


//Our main
int main(int argc, char **argv){
    
	if(argc != 3 || strcmp("-n",argv[1]) != 0){
		errorf("Not a valid format expected: [./multiplier -n NUM_BUFFERS]  \n");
		exit(EXIT_FAILURE);
	}

	NUM_BUFFERS = strtol(argv[2], NULL, 10);

	if(NUM_BUFFERS < 1){
		errorf("Not a valid number of buffers: the current number is %d , it should be greater than 0\n",NUM_BUFFERS);
		exit(EXIT_FAILURE);
	}

	if(NUM_BUFFERS < 12){
		warnf("There are [ %d ]  buffers, this might cause errors in execution time, it is recomended to use at least [ 12 ] buffers\n",NUM_BUFFERS);
	}	

	infof("The solution has started running with %d buffers\n",NUM_BUFFERS);

	buffers = (long **)malloc(NUM_BUFFERS * sizeof(long *));
	mutexes = (pthread_mutex_t *) malloc(NUM_BUFFERS * sizeof(pthread_mutex_t));

	for (int i = 0; i < NUM_BUFFERS; i++) {
		pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
		mutexes[i] = mutex;
		pthread_mutex_init(&mutexes[i], NULL);
	}
	long *matrixA, *matrixB;
	matrixA = readMatrix("matA.dat");
	matrixB = readMatrix("matB.dat");
	result= multiply(matrixA, matrixB);
	saveResultMatrix(result);
	infof("The resulting matrix was saved in 'result.dat' file\n");

	free(matrixA);
	free(matrixB);
	free(buffers);
	free(mutexes);
	free(result);

    return 0;
}
