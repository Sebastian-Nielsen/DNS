package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type SafeSet struct {
	mu sync.Mutex
	values map[net.Conn]bool
}
func (s *SafeSet) add(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[conn] = true
}
func (s *SafeSet) delete(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.values, conn) // delete conn from the set of openConnectitons
}

var outbound = make(chan string)
var openConnections = SafeSet{}

func listen(conn net.Conn) {
	defer conn.Close()
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if (err != nil) {
			fmt.Println("Error: " + err.Error())
			return
		} else {
			// No error, send the message to the channel
			fmt.Print("From Client:", string(msg))
			titlemsg := strings.Title(msg)

			outbound <- titlemsg
			//conn.Write([]byte(titlemsg))
		}
	}
}

func broadcast() {  // Whenever there is a msg in the outbounds channel send it
	i := 0
	for {
		msg := <- outbound
		fmt.Println("Broadcasting msg: '" + msg + "'")
		conn.Write([]byte(msg))
		i++
	}
}

func main() {  // SERVER

	fmt.Println("Listening for connection...")
	//addrs, _ := net.LookupHost("localhost")
	//fmt.Println(addrs)

	ln, _ := net.Listen("tcp", ":18081")
	defer ln.Close()

	go broadcast()

	// Continously listen for new connection requests
	for {
		conn, _ := ln.Accept()
		openConnections.add(conn); // Add the connection to the set of open connections
		fmt.Println("Got a new connection...")
		go listen(conn)
	}
}
