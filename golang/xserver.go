package main

import (
    "fmt"
    "net"
    "encoding/gob"
    "strconv"
)

type P struct {
    number int
    O string
}

func handleConnection(conn net.Conn) {
    dec := gob.NewDecoder(conn)
    defer conn.Close()
    for {
    	p := &P{}
    	err := dec.Decode(p)
	if (err != nil) {
	   return
	}
    	fmt.Printf(strconv.Itoa(p.number))
    }
}

func main() {
    fmt.Println("start");
    ln, _ := net.Listen("tcp", ":8080")
    for {
        conn, _ := ln.Accept()
        go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
    }
}