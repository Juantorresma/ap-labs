void snake_place(snake *a_snake, snake_pos *pos1)
{
    int i;
    for (i=0; i<a_snake->size; ++i)
    {
        gotoxy(a_snake->tail_X,a_snake->tail_Y);
        printf("%c",a_snake->symbol);
        pos1->X[i]=a_snake->tail_X;
        pos1->Y[i]=a_snake->tail_Y;
        a_snake->tail_X+=1;
    }

}

void snake_move(snake *a_snake, snake_pos *pos1, int *score){

        move_head(a_snake,pos1);
        if (!(state[(a_snake->head_X)][(a_snake->head_Y)]=='F')){
            move_tail(a_snake,pos1);
        } else {
            a_snake->size++;
            *score=*score+1;
            state[(a_snake->head_X)][(a_snake->head_Y)]=' ';
        }
}

int game_over(snake *a_snake, snake_pos *pos1){
    int i;
    for (i=0; i<a_snake->size-1; ++i){
        if ((pos1->X[i]==a_snake->head_X) && (pos1->Y[i]==a_snake->head_Y))
            return 1;
    }
    if ((a_snake->head_X==NCOLS) || (a_snake->head_X==1) || (a_snake->head_Y==NROWS) || (a_snake->head_Y==1))
        return 1;
    return 0;
}