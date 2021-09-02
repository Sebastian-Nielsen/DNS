package main
 
import ( "net" ; "fmt" ; "bufio" )

func main() {
     addrs, _ := net.LookupHost("www.google.com")
     addr := addrs[0]
     fmt.Println(addr)
     conn, err := net.Dial("tcp", addr+":80")
     if (conn!=nil) { defer conn.Close() }
     if (err!=nil) { panic(0) }
     fmt.Fprintf(conn, "GET /search?q=Secure+Distributed+Systems HTTP/1.1\n")
     fmt.Fprintf(conn, "HOST: www.google.com\n")
     fmt.Fprintf(conn, "\n")
     for {
     	msg, err := bufio.NewReader(conn).ReadString('\n')
     	if err != nil { panic(1) }
	fmt.Println(msg)
     }
}
