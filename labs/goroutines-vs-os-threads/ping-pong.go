package main

import (
	"fmt"
	"time"
	"flag"
	"strconv"
	"os"
)

func main() {
	  var pings string;
  	ping := make(chan int)
	  pong := make(chan int)
  	done := make(chan struct{})
    flag.Parse()
    pings = flag.Arg(0)
    pingpongCount, err := strconv.Atoi(pings)
    
    if err != nil {	
		    fmt.Println("Instructions: go run ping-pong.go <number of pings>")
        os.Exit(2)
    }

	startTime := time.Now()

	go func() {
  
		for n := 0; n < pingpongCount; n++ {
			ping <- n
			<-pong
		}
    
		close(ping)
		close(done)
	}()

	go func() {
  
		for n := range ping {
			//fmt.Printf("Pong!\n")
			pong <- n
		}

		close(pong)
	}()

	<-done
	endTime := time.Now()
	deltaT := endTime.Sub(startTime)
	time := float64(deltaT.Nanoseconds()) / 1000000000.0
	rate := float64(pingpongCount) / time
	fmt.Printf("Time: %v \t Message number: %v Replies: \t%f \n", deltaT, pingpongCount, rate)
}
