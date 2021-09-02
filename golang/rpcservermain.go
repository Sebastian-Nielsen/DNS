package main

import (
     "net" ;
     "net/rpc" ;
     "log" ;
     "net/http" ;
     "fmt"

     "./rpcserver" 
)

func main() {


     arith := new(rpcserver.Arith)
     server := rpc.NewServer()
     server.RegisterName("Arith", arith)
     
     rpc.HandleHTTP()


     l, e := net.Listen("tcp", ":1234")
     if e != nil {
     	fmt.Printf("e:",e)
     	log.Fatal("listen error:", e)
     }
     go http.Serve(l, nil)


     fmt.Println("Server up!")

     client, err := rpc.DialHTTP("tcp", "127.0.0.0" + ":1234")

     fmt.Println("Server up!")

     if err != nil {
     	log.Fatal("dialing:", err)
     }
     defer client.Close()
     fmt.Println("Done!")
}