package main
 
import ( "net" ; "fmt" )

func main() {
  ln, _ := net.Listen("tcp", ":")
  defer ln.Close()
  _, port, _ := net.SplitHostPort(ln.Addr().String())
  fmt.Println("Listening on port " + port)
}
