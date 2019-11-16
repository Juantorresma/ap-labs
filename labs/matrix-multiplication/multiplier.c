#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <strings.h>
#include "logger.h"

long * readMatrix(char *filename);
long * getColumn(int col, long *matrix);
long * getRow(int row, long *matrix);
int getLock();
int releaseLock(int lock);
long dotProduct(long *vec1, long *vec2);
long * multiply(long *matA, long *matB);
int saveResultMatrix(long *result);
void *threadFunc(void *arg);

#define MATRIX_DIMENSION 2000
#define MATRIX_NUM_ELEMENTS 4000000

int NUM_BUFFERS;
long **buffers;
pthread_mutex_t *mutexes;
long *resMul;
pthread_t threads[MATRIX_DIMENSION];

struct vectorStruct {
	int rowPos;
	int colPos;
	long *res;
    long *matA;
	long *matB;
};

long * readMatrix(char *filename){
    
    FILE *f;

    if ((f = fopen(filename,"r")) == NULL){
        errorf("Error! opening file");
        // Program exits if the file pointer returns NULL.
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
		warnf("col is an invalid number, must be from 0 to 2000\n");
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
		warnf("row is an invalid number, must be from 0 to 2000\n");
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

    long *resMul = (long *)malloc(MATRIX_NUM_ELEMENTS * sizeof(long));
	
    int i, 
        j;

	for ( i = 0; i < MATRIX_DIMENSION; i++) {
		for ( j = 0; j < MATRIX_DIMENSION; j++) {
			struct vectorStruct *vector;
            vector =(struct vectorStruct *) malloc(sizeof(struct vectorStruct));

			vector->matA = matA;
			vector->matB = matB;
			vector->rowPos = i + 1;
			vector->colPos = j + 1;
			vector->res = resMul;

			pthread_create(&threads[j], NULL , threadFunc , (void *)vector);
		}
//		printf("%d\n",i);

		for ( j = 0; j < MATRIX_DIMENSION; j++)
			pthread_join(threads[j], NULL);

		fflush(stdout);
	}

	return resMul;
}

int saveResultMatrix(long *result){

    FILE *f = fopen("result.dat", "w");
	if (f == NULL) {
		errorf("error writing 'result.dat' file\n");
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

	buffers[lockRow] = getRow(currVec->rowPos, currVec->matA);
	buffers[lockCol] = getColumn(currVec->colPos, currVec->matB);
	
	i = ((((currVec->rowPos - 1) * MATRIX_DIMENSION) + currVec->colPos) - 1);
	currVec->res[i] = dotProduct(buffers[lockRow], buffers[lockCol]);

	free(buffers[lockRow]);
	free(buffers[lockCol]);
	free(arg);

	while (releaseLock(lockRow) != 0);
	while (releaseLock(lockCol) != 0);

	return NULL;
}

int main(int argc, char **argv){
    
	if(argc != 3 || strcmp("-n",argv[1]) != 0){
		errorf("Invalid format [./multiplier -n NUM_BUFFERS] expected  \n");
		exit(EXIT_FAILURE);
	}

	NUM_BUFFERS = strtol(argv[2], NULL, 10);

	if(NUM_BUFFERS < 1){
		errorf("Invalid number of buffers : %d , must be greater than 0\n",NUM_BUFFERS);
		exit(EXIT_FAILURE);
	}

	if(NUM_BUFFERS < 12){
		warnf("[ %d ] number of buffers might cause errors in execution time, it is recomended to use (minimum) [ 12 ] buffers\n",NUM_BUFFERS);
	}	

	infof("Program start running with [ %d ] number of buffers\n",NUM_BUFFERS);

	buffers = (long **)malloc(NUM_BUFFERS * sizeof(long *));
	mutexes = (pthread_mutex_t *) malloc(NUM_BUFFERS * sizeof(pthread_mutex_t));

	for (int i = 0; i < NUM_BUFFERS; i++) {
		pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
		mutexes[i] = mutex;
		pthread_mutex_init(&mutexes[i], NULL);
	}

	infof("Threads initialized correctly \n");

	long *mA, *mB;
	mA = readMatrix("matA.dat");
	mB = readMatrix("matB.dat");

	infof("Matrices read correctly\n");

	infof("Computing multiplication\n");
	resMul= multiply(mA, mB);
	infof("Multiplication succesfully computed\n");

	saveResultMatrix(resMul);
	infof("Matrix saved in 'result.dat' file\n");

	free(mA);
	free(mB);
	free(buffers);
	free(mutexes);
	free(resMul);

    return 0;
}