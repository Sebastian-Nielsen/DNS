package main

import ( "fmt" ; "account" )

func main() {
     l := account.MakeLedger()
     t := account.SignedTransaction{ To: "Jesper", From: "Ivan", Amount: 1000000 }
     l.SignedTransaction(&t) 
     fmt.Println(l)
}
