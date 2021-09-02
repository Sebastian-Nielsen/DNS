package main

import (
    "fmt"
    "net"
    "encoding/gob"
    "bufio"
    "os"
)

type P struct {
    number int
    O string
}

func main() {
    fmt.Println("start client");
    conn, _ := net.Dial("tcp", "localhost:8080")
    defer conn.Close()
    encoder := gob.NewEncoder(conn)
    reader := bufio.NewReader(os.Stdin)
    var err error

    for i:=0;; i++ {
        fmt.Print("> ")
    	p := P{}
        p.O, err = reader.ReadString('\n')
	p.number = i
        if err!=nil || p.O == "quit\n" { return }
	fmt.Println(p)
    	encoder.Encode(&p)
    }
    fmt.Println("done");
}
