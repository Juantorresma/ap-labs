package main

import (
	"./snake"
	"os"
	"fmt"
	"strconv"
	)

func main() {
	if len(os.Args)==3{
		if fN,err := strconv.Atoi(os.Args[1]); err!=nil{
			fmt.Println(err)
		} else{
			snake.NewGame(fN).Start()
		}

	} else{
		fmt.Println("Specify the food quantity and the enemy quantity")
		fmt.Println("./main #F #E")
	}

}
