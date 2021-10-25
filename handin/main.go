package main

import (
	. "DNO/handin/Helper"
	. "DNO/handin/Peernode"
	"net"
)


func main() {
	var peerNode = PeerNode{
		OpenConnections:     SafeSet_Conn{ Values: make(map[net.Conn]bool) },
		PeersInArrivalOrder: SafeArray_string{},
		MessagesSent:        SafeSet_string{ Values: make(map[string  ]bool) },
		Ipc:                 IPC{ ConnToEncDecPair: make(map[net.Conn]EncoderDecoderPair) },
		TestMock:            Mock{ ShouldMockInput: false, ShouldPrintDebug: true },
	}

	peerNode.Start("", "")

	println("---------------------\nTerminating...")
}

