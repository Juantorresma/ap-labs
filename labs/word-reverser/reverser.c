//A01227885

#include <stdio.h>

//here i generate  the new word recursively
void reverse(char *a, int min, int max) {
    if(max-min < 1){return;}
    else {
    change(&a[min], &a[max]);
    reverse(a, min+1, max-1);
    }
}

//here i chage the corresponding letters
void change(char *a, char *b){
    char tmp = *a;
    *a = *b;
    *b = tmp;
}

//as my word is now an array, i made this func to print it
void printArray(char *a, int length){
    for(int i = 0; i < length; i++){
    printf("%c", a[i]);
    }
}


//here is my magic
int main(){
    char input[500];
    char c;
    printf("Enter a word: \n");
    for(int i = 0, c = getchar(); c != EOF; i++){
        input[i] = c;
        c =  getchar();
        if(c == '\n'){
            char word[i+1];
            for(int j = 0; j < i+1; j++){
                word[j] = input[j];
            }
            int length = sizeof(word)/sizeof(char);
            reverse(word, 0, length-1);
            printf("The reversed word is: \n");
            printArray(word, length);
            printf("\n");
            printf("Please enter a new word: \n");
            i = -1;
            continue;
        }
    }
    return 0;
}


