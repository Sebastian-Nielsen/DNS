package main

import ( "fmt" ; "strconv" )

var sendernames = [5]string{"Alice", "Bob", "Chloe", "Dennis", "Elisa"}

var receivernames = [5]string{"Frederik", "Gary", "Hailey", "Isabel", "Jesper"}

func send(c chan string, myname string) {
     for i:=0; i<1000; i++ {
     	 // you send on a channel using <-
     	 c<- myname + "#" + strconv.Itoa(i)
     }
}

func receive(c chan string, myname string) {
     i:=0
     for {
     	 // you also receive from channel using <-
     	 msg := <-c
	 fmt.Println(myname + "#" + strconv.Itoa(i) + " " + msg)
	 i++
     }
}

func main() {
  c := make(chan string)
  for i:=0; i<5; i++ {
      go send(c, sendernames[i])
      go receive(c, receivernames[i])
  }
  // we do one blocking call to avoid that the program terminates
  receive(c, "Kacey")
}
