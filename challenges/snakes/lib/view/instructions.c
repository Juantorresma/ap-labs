
void print_instructions(){
    char e;
    system("clear");
    printf("\t\t\t\t     M U L T I    S N A K E\n\n\n");
    printf("\t\t\t\t\tInstructions: \n");
    printf("Use the keys W, A, S, D to move the main snake. Eat the food and don't touch any other snake. \n");
    printf("Press Ctrl + C to end the game \n\n");
    printf("\t\tEnter the number of fruits: ");
    scanf("%d", &num_fruits);
    printf("\t\tEnter the number of enemies: ");
    scanf("%d", &num_enemies);
    printf("\t\tPress enter to continue \n\n\nCredits: Karol Gutierrez and Enrique Bovio\n\n");
    scanf("%c", &e);
}
