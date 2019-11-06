//pipeline excersice a01227885

package main

import(
	"time"
	"fmt"
)

const layout = "15:04:05.000000"

func main() {
	numPipes := 1000 
	var channels = make([]chan string, 0)
  
	for i := 0; i < numPipes; i++ {
		channels = append(channels, make(chan string))
	}
  
	go firstPipe(channels[0])
  
	for i := 1; i < numPipes; i++ {
		go middlePipe(channels[i-1], channels[i])
	}
  
	startTime, _ := time.Parse(layout, <-channels[numPipes-1])
	diff := time.Since(startTime)
	fmt.Println("Time it took for our message to pass through all channels and goroutines was:", diff)
  
}

func firstPipe(chOut chan string) {
	chOut <- time.Now().Format(layout)
	close(chOut)
}

func middlePipe(chIn chan string, chOut chan string) {
	t := <-chIn
	chOut <- t
	close(chOut)
}
