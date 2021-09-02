package main

import ( "fmt" ; "account" )

func main() {
     l := account.MakeLedger()
     t := account.Transaction{ To: "Jesper", From: "Ivan", Amount: 1000000}
     l.Transaction(&t) 
     fmt.Println(l)
}
