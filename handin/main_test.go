package main

import (
	. "DNO/handin/Helper"
	"fmt"
	"net"
	"testing"
	"time"
)

func createPeerNode( shouldMockInput bool) PeerNode {
	return PeerNode{
		OpenConnections: SafeSet_Conn{   Values: make(map[net.Conn]bool) },
		MessagesSent:    SafeSet_string{ Values: make(map[string  ]bool) },
		TestMock:        Mock{ ShouldMockInput: shouldMockInput },
	}
}

const peerNode1_port = "50001"
const peerNode2_port = "50002"


func TestPeer1ReceivesMsgFromPeer2(t *testing.T) {

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)

	peerNode1.TestMock.SimulatedInputString = "no_port"
	peerNode2.TestMock.SimulatedInputString = peerNode1_port

	go peerNode1.Start(peerNode1_port)

	time.Sleep(3 * time.Second)

	go peerNode2.Start(peerNode2_port)

	time.Sleep(3 * time.Second)
	fmt.Println(peerNode1.OpenConnections.Values, len(peerNode1.OpenConnections.Values))
	if  len(peerNode1.OpenConnections.Values) != 1 &&
		len(peerNode2.OpenConnections.Values) != 1 {
		t.Errorf("asdflkjasdflk sadfkl asdf")
	}
}

func TestPeer1CanConnectToPeer2(t *testing.T) {


	peerNode1.TestMock.SimulatedInputString = "dont_connect_to_any_peer"
	peerNode2.TestMock.SimulatedInputString = peerNode1_port

	go peerNode1.Start(peerNode1_port)

	time.Sleep(3 * time.Second)

	go peerNode2.Start(peerNode2_port)

	time.Sleep(3 * time.Second)
	fmt.Println(peerNode1.OpenConnections.Values, len(peerNode1.OpenConnections.Values))
	if  len(peerNode1.OpenConnections.Values) != 1 &&
		len(peerNode2.OpenConnections.Values) != 1 {
		t.Errorf("asdflkjasdflk sadfkl asdf")
	}
}









