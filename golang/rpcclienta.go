package main

import ( "log" ; "net/rpc" ; "bufio" ; "os" ; "fmt")

func PrintWhenReady(call *rpc.Call) {
   <-call.Done // this channel unblocks when the call returns
   if call.Error != nil { log.Fatal(call.Error) }
   fmt.Println("Strings printed at server: ", reply)
}

var reply int

func main() {
   client, err := rpc.Dial("tcp", "localhost:42587")	
   if err != nil { log.Fatal(err) }

   in := bufio.NewReader(os.Stdin)
   for {
      line, _, _ := in.ReadLine()		

      // Asynchronous call
      call := client.Go("PrintAndCount.GetLine", line, &reply, nil)
      go PrintWhenReady(call) // Handles the reply when ready

      fmt.Println("See, I can still do stuff!")
      fmt.Println("See, I can still do stuff!")
      fmt.Println("OK, I'm bored...")
   }
}