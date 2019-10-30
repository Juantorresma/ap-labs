// Netcat2 is a read-only TCP client with channels
//A01227885 used my clock 2 as base
package main

//imports
import (
	"log"
	"net"
	"os"
	"strings"
)

//my main with locations  i will be using
func main() {

	var newYork, tokyo, london string
  
	for _, s := range os.Args[1:] {
		list := strings.Split(s, "=")
    
	    if list[0] == "NewYork" {
			newYork = list[1]
	    }else if list[0] == "Tokyo" {
			tokyo = list[1]
		}else if list[0] == "London" {
			london = list[1]
		}else{
			log.Fatal("Location is not valid ")
		}
    
	}
	var newYorkConnection, tokyoConnection, londonConnection net.Conn
	var errors error
	if newYork!=""{
		newYorkConnection, errors = net.Dial("tcp", newYork)
		if errors != nil {
			log.Fatal(errors)
		}
	}
	if tokyo!=""{
		tokyoConnection, errors = net.Dial("tcp", tokyo)
		if errors != nil {
			log.Fatal(errors)
		}
	}
	if london!=""{
		londonConnection, errors = net.Dial("tcp", london)
		if errors != nil {
			log.Fatal(errors)
		}
	}
	done := make(chan int)
	go func(newYorkConnection, tokyoConnection, londonConnection net.Conn) {
		for{
			buffer := make([]byte, 1400)
			if newYork!=""{
				dataSize, errors := newYorkConnection.Read(buffer)
				if errors != nil {
					log.Fatal(errors)
					break
				}
				data := buffer[:dataSize]
                log.Print(string(data))
			}
			if tokyo!=""{
				dataSize, errors := tokyoConnection.Read(buffer)
				if errors != nil {
					log.Fatal(errors)
					break
				}
				data := buffer[:dataSize]
                log.Print(string(data))
			}
			if london!=""{
				dataSize, errors := londonConnection.Read(buffer)
				if errors != nil {
					log.Fatal(errors)
					break
				}
				data := buffer[:dataSize]
                log.Print(string(data))
			}
		}
		log.Println("done")
		done <- 2 // signal the main goroutine
	}(newYorkConnection, tokyoConnection, londonConnection)

	x := 1
	x = <-done 
	log.Println("Connection Closed, value: ", x)
	close(done)
}
