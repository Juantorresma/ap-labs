typedef struct snake_pos{
    int Y[NROWS*NCOLS];
    int X[NROWS*NCOLS];

}snake_pos;

void pos_init(snake_pos *pos1){
    memset(pos1, 0, sizeof(*pos1));
}
