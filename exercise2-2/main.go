package main;

import (
	"fmt"
)


func main() {
	fmt.Println("Starting program ...\n---------------------")

	var ip string
	var port string
	fmt.Println("Please write an IP and port:")
    fmt.Scanln(&ip)
	fmt.Println("Please write the port number")
	fmt.Scanln(&port)

	fmt.Println("'ip:port' number is '" + ip + ":" + port + "'")
	fmt.Println("---------------------\nTerminating...")
}
