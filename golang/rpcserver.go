package main

import ( "fmt" ; "log" ; "net" ; "net/rpc" ; "os" ; "bufio" ; "time")

type PrintAndCount struct {
   x int
}

func (l *PrintAndCount) GetLine(line []byte, cnt *int) error {
   HardTask()   
   l.x++
   fmt.Println(string(line))
   *cnt = l.x
   return nil
}

func HardTask() {
   time.Sleep(5 * time.Second)           
}

func MakeTCPListener() *net.TCPListener {
   addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:42587")
   if err != nil { log.Fatal(err) }

   inbound, err := net.ListenTCP("tcp", addy)
   if err != nil { log.Fatal(err) }
   return inbound
}

func main() {
   // Register how we receive incoming connections
   go rpc.Accept(MakeTCPListener())

   // Register a PrintAndCount object 
   rpc.Register(new(PrintAndCount))

   // Avoid terminating
   fmt.Println("Press any key to terminate server")
   bufio.NewReader(os.Stdin).ReadLine()   
}