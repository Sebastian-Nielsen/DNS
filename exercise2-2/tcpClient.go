package main

import (
	"bufio";
	"fmt";
	"net";
	"os"
)

var conn net.Conn

func send(conn net.Conn) {
	// Continously ask for input to send to the server
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if text == "quit\n" { return }
		if err != nil { return }
		fmt.Fprintf(conn, text) // Send the inputted text to the connection
	}
}
func receive(conn net.Conn) {
	// Continously listen for messages from the server
	// Upon receipt of messages, display them
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil { return }
		fmt.Print("From server: " + msg)
	}
}

func main() {  // CLIENT
	conn, _ = net.Dial("tcp", "127.0.0.1:18081")
	defer conn.Close()

	send(conn)

	fmt.Println("Terminating ...")
}

