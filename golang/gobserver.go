package main

import ( "net" ; "fmt" ; "encoding/gob" ; "io" ; "log" )

type ToSend struct {
     Msg string // only exported variables are sent, so start the ...
     Number int // ... name of the fields you want send by a capital letter
}

func handleConnection(conn net.Conn) {
  defer conn.Close()
  msg := &ToSend{}
  dec := gob.NewDecoder(conn)
  for {
    err := dec.Decode(msg)
    if (err == io.EOF) {
      fmt.Println("Connection closed by " + conn.RemoteAddr().String())
      return
    }
    if (err != nil) {
      log.Println(err.Error())
      return
    }
    fmt.Println("From " + conn.RemoteAddr().String() + ":\n", msg)
  }
}

func main() {
  fmt.Println("Listening for connection...")
  ln, _ := net.Listen("tcp", ":18081")
  defer ln.Close()
  conn, _ := ln.Accept()
  fmt.Println("Got a connection from ", conn.RemoteAddr().String())
  handleConnection(conn)
}
