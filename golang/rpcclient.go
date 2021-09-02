package main

import ( "log" ; "net/rpc" ; "bufio" ; "os" ; "fmt")

func main() {
   client, err := rpc.Dial("tcp", "localhost:42587")	
   if err != nil { log.Fatal(err) }

   in := bufio.NewReader(os.Stdin)
   for {
      line, _, err := in.ReadLine()		
      var reply int


      // Synchronous call
      err = client.Call("PrintAndCount.GetLine", line, &reply)

      if err != nil { log.Fatal(err) }
      fmt.Println("Strings printed at server: ", reply)
   }
}