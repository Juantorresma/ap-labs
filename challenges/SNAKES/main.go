package main

import (
	"./snake"
	"os"
	"fmt"
	"strconv"
	)

func main() {
	if len(os.Args)==2{
		if fN,err := strconv.Atoi(os.Args[1]); err!=nil{
			fmt.Println(err)
		} else{
			snake.NewGame(fN).Start()
		}

	} else{
		fmt.Println("Specify the food quantity")
		fmt.Println("./main #F")
	}

}
