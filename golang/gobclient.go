package main

import ( "bufio" ; "fmt" ; "net" ; "os" ; "encoding/gob" )

type ToSend struct {
     Msg string // only exported variables are sent, so start the ...
     Number int // ... name of the fields you want sent by a capital letter
}

func main() {
     ts := &ToSend{}
     conn, _ := net.Dial("tcp", "127.0.0.1:18081")
     defer conn.Close()
     enc := gob.NewEncoder(conn)
     for i:=0;; i++ {
		fmt.Print("> ")
			reader := bufio.NewReader(os.Stdin)
		m, err := reader.ReadString('\n')
		if err!=nil || m == "quit\n" { return }
		ts.Msg = m
		ts.Number = i
		enc.Encode(ts)
     }
}
