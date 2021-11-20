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
type IPC struct {
	ConnToEncDecPair map[net.Conn]EncoderDecoderPair
}
type EncoderDecoderPair struct {
	Encoder *gob.Encoder
	Decoder *gob.Decoder
}

func (ipc *IPC) Send(packet Packet, conn net.Conn) bool {
	ok := true
	// fmt.Println("\t[IPC:Send]    before encoding:", packet)

	connHasAnEncoderDecoderPair := ipc.ConnToEncDecPair[conn] != (EncoderDecoderPair{})
	if !connHasAnEncoderDecoderPair {
		ipc.ConnToEncDecPair[conn] = EncoderDecoderPair{
			Encoder: gob.NewEncoder(conn),
			Decoder: gob.NewDecoder(conn),
		}
	}


	enc := ipc.ConnToEncDecPair[conn].Encoder
	err := enc.Encode( packet )
	if err != nil {
		fmt.Println("Ipc Send err:", err)
		ok = false
	}

	return ok
}
func (ipc *IPC) Receive(conn net.Conn) (Packet, bool) {
	var packet Packet
	ok := true

	connHasAnEncoderDecoderPair := ipc.ConnToEncDecPair[conn] != (EncoderDecoderPair{})
	if !connHasAnEncoderDecoderPair {
		ipc.ConnToEncDecPair[conn] = EncoderDecoderPair{
			Encoder: gob.NewEncoder(conn),
			Decoder: gob.NewDecoder(conn),
		}
	}

	dec := ipc.ConnToEncDecPair[conn].Decoder
	err := dec.Decode(&packet)
	//fmt.Println("\t1[IPC:Receive] Decoded packet:", err, packet)
	if err != nil {
		fmt.Println("Ipc receive err:", err)
		ok = false
	}
	//fmt.Println("\t2[IPC:Receive] Decoded packet:", err, packet)
	return packet, ok
}
