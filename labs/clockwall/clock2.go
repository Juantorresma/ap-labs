// Clock2 is a concurrent TCP server that periodically writes the time.
package main

//imports
import (
	"io"
	"log"
	"net"
	"time"
	"strconv"
	"os"
)

//modified handle con for errors
func handleConn(c net.Conn, tz string) {
	defer c.Close()
	loc, er := time.LoadLocation(tz)
	if er != nil {
		log.Fatal(er)
		return // e.g., client disconnected
	}
	for {
		_, err := io.WriteString(c, tz+": "+time.Now().In(loc).Format("15:04:05\n"))
		if err != nil {
			log.Fatal(err)
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

//my new main with a flag which i will be using, and the name
func main() {
	flag:= false
	var name string
	for _, s := range os.Args {
	    if s == "-port" {
	        flag=true
			continue
	    }else if flag{
			flag= false
			i, err := strconv.Atoi(s)
			if err != nil || i>65535 || i<0{
				log.Fatal(err)
			}
			name = "localhost:"+strconv.Itoa(i)
		}
	}
	tz := os.Getenv("TZ")
	listener, err := net.Listen("tcp", name)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}

		go handleConn(conn, tz) // handle connections concurrently
	}
}
