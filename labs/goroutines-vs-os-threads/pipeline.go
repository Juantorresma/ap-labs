//pipeline excersice a01227885

package main

import (
	"fmt"
	"time"
	"os"
	"strconv"
)

func main() {
	head := make(chan time.Time)
	last := head
    if err != nil {
        fmt.Println(err)
        return
    }
	for stageCount := 1; ; stageCount++ {
		go pipeLine(last, stageCount, f)
		head <- time.Now()

		temp := last
		last = make(chan time.Time)
		go connectPipes(temp, last)
	}
}

func pipeLine(last chan time.Time, stageCount int, writer *os.File) {
	startTime := <-last
	endTime := time.Now()
	fmt.Printf("Goroutine number: %d\t Time: %v\n", stageCount, endTime.Sub(startTime))
	writer.WriteString(strconv.Itoa(stageCount)+","+fmt.Sprint(endTime.Sub(startTime))+"\n")
}

func connectPipes(src, dst chan time.Time) {
	for t := range src {
		dst <- t
	}
}
