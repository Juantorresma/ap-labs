void set_borders(){
    int i;
    for (i=0; i<NROWS; ++i){
        gotoxy(0,i);
        printf("X");
        gotoxy(NCOLS,i);
        printf("X");
    }
    for (i=0; i<NCOLS; ++i){
        gotoxy(i,0);
        printf("X");
        gotoxy(i,NROWS);
        printf("X");
    }
}