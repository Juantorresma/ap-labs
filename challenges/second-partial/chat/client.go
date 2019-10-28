// A01227885 my modification for client
package main

//al the imports needed
import (
	"io"
	"log"
	"net"
	"os"
)

//mmy main where i will do everything
func main() {
	
	//firts i set my  string variables
	var svr, user string
	//then two booleands for my flags
	var svr_flag, user_flag bool
	for _, s := range os.Args {
	    if s == "-user" {
	        user_flag=true
			continue
	    }else if user_flag{
			user_flag= false
			user = s
	    }if s == "-server" {
		    	svr_flag=true
			continue
	    }else if svr_flag{
			svr_flag= false
			svr = s
		}
	}
	//check if they are not blanc
	if svr == "" || user == ""{
		log.Fatal("There is no server or user")
	}

	conn, err := net.Dial("tcp", svr)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	_, erro := conn.Write([]byte(user))
	if erro != nil {
		log.Fatal(erro)
	}

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
