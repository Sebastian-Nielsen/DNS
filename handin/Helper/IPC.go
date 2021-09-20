package Helper

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

/*
	Inter-process communication (IPC)

	Handles serialization/deserialization and sending of data, as well as error handling
*/
type IPC struct {
	mu     sync.Mutex
}

func (ipc* IPC) Send(packet Packet, conn net.Conn) bool {
	ipc.mu.Lock()
	defer ipc.mu.Unlock()

	ok := true
	//fmt.Println("\t[IPC:Send]    before encoding:", packet)
	enc := gob.NewEncoder(conn)
	err := enc.Encode( packet )
	if err != nil {
		fmt.Println("Ipc Send err:", err)
		ok = false
	}

	// A thread can only send a message every 0.5 seconds
	time.Sleep(40 * time.Millisecond)
	return ok
}
func (ipc* IPC) Receive(conn net.Conn) (Packet, bool) {
	var packet Packet
	ok := true
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&packet)
	//fmt.Println("\t[IPC:Receive] Decoded packet:", packet)
	if err != nil {
		fmt.Println("Ipc receive err:", err)
		ok = false
	}
	return packet, ok
}
