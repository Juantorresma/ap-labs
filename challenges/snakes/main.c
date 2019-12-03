#include <stdio.h>
#include <string.h>
#include <termios.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <ncurses.h>
#include <time.h>
#include <omp.h>

int num_enemies;
int num_fruits;
int i;
char **state;
int nThreads;

#include "lib/controller/config.c"
#include "lib/controller/print_to_console.c"
#include "lib/view/print_field.c"
#include "lib/view/instructions.c"
#include "lib/controller/keyboard_input.c"
#include "lib/model/food.c"
#include "lib/model/snake.c"
#include "lib/model/snake_positions.c"
#include "lib/controller/move_snake.c"
#include "lib/controller/game_step.c"


int main(){
    print_instructions();
      int snake_speed=100000;
      int score=0;
      state = malloc(NROWS*sizeof(char*));
      for(i=0; i<NCOLS; i++){
             state[i]=malloc(NCOLS*sizeof(char));
      }
      snake a_snake;
      snake_pos pos1;


      food *foods=malloc(num_fruits*sizeof(food));
      snake_init(&a_snake, 'O');
      pos_init(&pos1);
      for(i=0; i<num_fruits; i++){
             food_init(&foods[i]);
             state[foods[i].X][foods[i].Y]='F';
      }

      system("clear");
      system("stty -echo");
      snake_place(&a_snake,&pos1);
      set_borders();
      for(i=0; i<num_fruits; i++){
             food_print(&foods[i]);
      }

      system ("/bin/stty raw");

       omp_set_num_threads(num_enemies);
       #pragma omp parallel
       {
              snake enemy;
              snake_pos eneypos;
              int id = omp_get_thread_num();
              printf("Thread: %d\n", id);
       }

      while(!(game_over(&a_snake,&pos1))){
          while (!kbhit()){
                 usleep(snake_speed);
                 snake_move(&a_snake,&pos1,&score);
                 if (game_over(&a_snake,&pos1)) break;
          }
          a_snake.prev_direction=a_snake.direction;
          int w = getchar();
          if(w == UP || w == DOWN || w == LEFT || w == RIGHT)
            a_snake.direction=w;
     }
      system ("/bin/stty cooked");
      system("stty echo");
      system("clear");
  //    nThreads = omp_get_num_threads();
    //  printf("Number of threads = %d\n", nThreads);
      printf("\n\n Game finished, your score: %d \n\n", score);
      return 0;
}
