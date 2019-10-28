// A01227885 my server side of the chat with the base you gave us
package main

//al the imports
import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	"strconv"
	"os"
	"strings"
)

//given functions
type client chan<- string // an outgoing message channel

var (
	entering = make(chan user)
	leaving  = make(chan user)
	messages = make(chan mssg) // all incoming client messages
)

//a struct for the users whit user chanel and conections
type user struct {
usr string
chan  client
connect net.Conn
}

type mssg struct {
message string
chan  client
usr string
}


func broadcaster() {
	clients := make(map[string]user) // all connected clients
	badhosts := make(map[string]bool)
	missingAdmin := true
	var admin user
	for {
		select {
		case fisrtMMsg := <-messages:
			s := strings.Split(fisrtMMsg.message, " ")
			cmd := s[1]
			if cmd=="/users"{
				keys := make([]string, 0, len(clients))
				for key := range clients {
					keys = append(keys, key)
  				}
				fisrtMMsg.chan <- strings.Join(keys, ", ")
			}else if cmd=="/msg" {
				if len(s)>3{
					if to, found := clients[s[2]]; found {
						to.chan <- "@" + fisrtMMsg.usr + ": " + strings.Join(s[3:], " ")
					}else{
						fisrtMMsg.chan <- "irc-server > User isn't connected"
					}
				}else{
					fisrtMMsg.chan <- "irc-server > Usage: /msg <user> <msg>"
				}
			}else if cmd=="/time"{
				loc, e := time.LoadLocation("America/Mexico_City")
				if e != nil {
					log.Fatal(e)
				}
				fisrtMMsg.chan <- "irc-server > Local Time: America/Mexico_City " + time.Now().In(loc).Format("15:04:05")
			} else if cmd=="/user"{
				if len(s)==3{
					if to, found := clients[s[2]]; found {
						fisrtMMsg.chan <- "irc-server > username: " + to.usr + ", IP: " + strings.Split(to.connect.RemoteAddr().String(), ":")[0]
					}else{
						fisrtMMsg.chan <- "irc-server > User isn't connected"
					}
				}else{
					fisrtMMsg.chan <- "irc-server > Usage: /user <user>"
				}
			}else if cmd=="/kick"{
				if fisrtMMsg.usr == admin.use{
					if len(s)==3{
						if to, found := clients[s[2]]; found {
							log.Printf("[%s] was kicked", to.usr)
							to.chan <- "irc-server > You're kicked from this channel"
							to.chan <- "irc-server > Bad language is not allowed on this channel"
							to.connect.Close()
							for _, client := range clients {
								client.chan <- "irc-server > [" + to.usr + "] was kicked from channel for bad language policy violation"
							}
						}else{
							fisrtMMsg.chan <- "irc-server > User isn't connected"
						}
					}else{
						fisrtMMsg.chan <- "irc-server > Usage: /kick <user>"
					}
				}else{
					fisrtMMsg.chan <- "irc-server > Only admin can kick user"
				}
			}else{
			// Broadcast incoming message to all
			// clients' outgoing message channels.
				for _, client := range clients {
					if fisrtMMsg.usr!=client.usr{
						client.chan <- fisrtMMsg.message
					}
				}
			}

		case cli := <-entering:
			if _,found := clients[cli.usr]; found {
				log.Printf("Client trying to use occupied username: [%s|%s]", cli.usr, cli.connect.RemoteAddr().String())
				cli.chan <- "irc-server > This username is already in use, please try another one."
				badhosts[cli.connect.RemoteAddr().String()] = true
				cli.connect.Close()
			}else{
				log.Printf("New connected user [%s]", cli.usr)
				cli.chan <- "irc-server > Your user [" + cli.usr + "] is successfully logged"
				fisrtMMsg := "irc-server > " + cli.usr + " has arrived"
				for _, client := range clients {
					client.chan <- fisrtMMsg
				}
				clients[cli.usr] = cli
			}
			if missingAdmin{
				missingAdmin = false
				admin = cli
				admin.chan <- "irc-server > Congrats, you were the first user."
				admin.chan <- "irc-server > You're the new IRC Server ADMIN"
				log.Printf("[%s] was promoted as the channel ADMIN", admin.usr)
			}

		case cli := <-leaving:
			if _, found := badhosts[cli.connect.RemoteAddr().String()]; !found {
				delete(clients, cli.usr)
				log.Printf("[%s] left", cli.usr)
				fisrtMMsg:= "irc-server > " + cli.usr + " has left"
				for _, client := range clients {
					client.chan <- fisrtMMsg
				}
			}
			close(cli.chan)
			cli.connect.Close()
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	buffer := make([]byte, 1400)
	dataSize, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	who := string(buffer[:dataSize])

	ch <- "irc-server > Welcome to the Simple IRC Server"
	entering <- user{usr: who, chan: ch, connect: conn}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- mssg{who + ": " + input.Text(), ch, who}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- user{usr: who, chan: ch, connect: conn}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}


//my main function
func main() {
	//set flags and start
	portflag:= false
	ipflag:= false
	var ip, port, host string
	for _, s := range os.Args {
	    if s == "-port" {
	        portflag=true
			continue
	    }else if portflag{
			portflag= false
			i, err := strconv.Atoi(s)
			if err != nil || i>65535 || i<0{
				log.Fatal(err)
			}
			port = s
		}
		if s == "-host" {
	        ipflag=true
			continue
	    }else if ipflag{
			ipflag= false
			ip = s
		}
	}
	if ip!="" && port!=""{
		host = ip+":"+port
	}else{
		log.Fatal("No port or host")
		return
	}
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started at ", host)
	log.Println("Ready for receiving new clients")
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
