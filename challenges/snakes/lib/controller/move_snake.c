void move_head(snake *snake1, snake_pos *pos1){
    switch (snake1->direction){
            case RIGHT:
                if (snake1->prev_direction==LEFT){
                    snake1->head_X--;
                    break;
                }
                    snake1->head_X++;
                    break;
            case LEFT:
                if (snake1->prev_direction==RIGHT){
                    snake1->head_X++;
                    break;
                }
                    snake1->head_X--;
                    break;
            case UP:
                if (snake1->prev_direction==DOWN){
                    snake1->head_Y++;
                    break;
                }
                    snake1->head_Y--;
                    break;
            case DOWN:
                if (snake1->prev_direction==UP){
                    snake1->head_Y--;
                    break;
                }
                snake1->head_Y++;
                break;
            default:
                 break;
        }
        pos1->X[snake1->size]=snake1->head_X;
        pos1->Y[snake1->size]=snake1->head_Y;
        gotoxy(pos1->X[snake1->size],pos1->Y[snake1->size]);
        printf("%c",snake1->symbol);
}

void move_tail(snake *snake1, snake_pos *pos1){
    int i;
    gotoxy(pos1->X[0],pos1->Y[0]);
    printf(" ");
    for (i=0; i<snake1->size; ++i){
        pos1->X[i]=pos1->X[i+1];
        pos1->Y[i]=pos1->Y[i+1];
    }
}