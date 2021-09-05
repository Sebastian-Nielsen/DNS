package Helper

import (
	"encoding/gob"
	"fmt"
	"net"
)

/*
	Inter-process communication (IPC)

	Handles serialization/deserialization and sending of data, as well as error handling
*/
type IPC struct {}

func (ipc* IPC) Send(packet Packet, conn net.Conn) bool {
	ok := true
	enc := gob.NewEncoder(conn)
	err := enc.Encode( packet )
	if err != nil {
		fmt.Println("Ipc Send err:", err)
		ok = false
	}
	return ok
}
func (ipc* IPC) Receive(conn net.Conn) (Packet, bool) {
	var packet Packet
	ok := true
	enc := gob.NewDecoder(conn)
	err := enc.Decode(&packet)
	if err != nil {
		fmt.Println("Ipc receive err:", err)
		ok = false
	}
	return packet, ok
}
