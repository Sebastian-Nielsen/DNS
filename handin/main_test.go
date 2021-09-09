package main

import (
	. "DNO/handin/Helper"
	"net"
	"testing"
	"time"
)

/*
	debug_mode has to be set to false for tests to pass for some mystical reason
 */

func createPeerNode( shouldMockInput bool) PeerNode {
	return PeerNode{
		OpenConnections: SafeSet_Conn{   Values: make(map[net.Conn]bool) },
		MessagesSent:    SafeSet_string{ Values: make(map[string  ]bool) },
		TestMock:        Mock{ ShouldMockInput: shouldMockInput },
	}
}



func TestMessageIsSentFromNode1To3Via2(t *testing.T) {

	const peerNode1_port = "50007"
	const peerNode2_port = "50008"
	const peerNode3_port = "50009"

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)
	peerNode3 := createPeerNode(true)


	peerNode1.TestMock.SimulatedInputString = "no_port"
	peerNode2.TestMock.SimulatedInputString = peerNode1_port
	peerNode3.TestMock.SimulatedInputString = peerNode2_port

	go peerNode1.Start(peerNode1_port)
	time.Sleep(500 * time.Millisecond)
	go peerNode2.Start(peerNode2_port)
	time.Sleep(500 * time.Millisecond)
	go peerNode3.Start(peerNode3_port)
	time.Sleep(500 * time.Millisecond)

	peerNode1.TestMock.SimulatedInputString = "some_msg"
	time.Sleep(500 * time.Millisecond)

	if !peerNode3.MessagesSent.Contains("some_msg") {
		t.Errorf("peerNode3 didn't receive peerNode1's msgs")
	}
}




func TestLatercomerNodeEventuallyGetsAllMsgs(t *testing.T) {

	const peerNode1_port = "50005"
	const peerNode2_port = "50006"

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)

	peerNode1.TestMock.SimulatedInputString = "no_port"
	peerNode2.TestMock.SimulatedInputString = peerNode1_port

	go peerNode1.Start(peerNode1_port)


	time.Sleep(1 * time.Second)
	peerNode1.TestMock.SimulatedInputString = "some_msg_1"
	time.Sleep(1 * time.Second)
	peerNode1.TestMock.SimulatedInputString = "some_msg_2"
	time.Sleep(1 * time.Second)
	peerNode1.TestMock.SimulatedInputString = "some_msg_3"
	time.Sleep(1 * time.Second)

	go peerNode2.Start(peerNode2_port)

	time.Sleep(1 * time.Second)

	if  !peerNode2.MessagesSent.Contains("some_msg_1") ||
		!peerNode2.MessagesSent.Contains("some_msg_2") ||
		!peerNode2.MessagesSent.Contains("some_msg_3") {
		t.Errorf("peerNode2 didn't receive peerNode1's msgs")
	}
}


func TestPeer1ReceivesMsgFromPeer2(t *testing.T) {

	const peerNode1_port = "50001"
	const peerNode2_port = "50002"

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)

	peerNode1.TestMock.SimulatedInputString = "no_port"
	peerNode2.TestMock.SimulatedInputString = peerNode1_port

	go peerNode1.Start(peerNode1_port)
	go peerNode2.Start(peerNode2_port)

	time.Sleep(1 * time.Second)

	peerNode1.TestMock.SimulatedInputString = "some_msg"

	time.Sleep(1 * time.Second)

	if !peerNode2.MessagesSent.Contains("some_msg") {
		t.Errorf("peerNode2 didn't receive peerNode1's msg")
	}

}

func TestPeer1CanConnectToPeer2(t *testing.T) {

	const peerNode1_port = "50003"
	const peerNode2_port = "50004"

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)

	peerNode1.TestMock.SimulatedInputString = "dont_connect_to_any_peer"
	peerNode2.TestMock.SimulatedInputString = peerNode1_port

	go peerNode1.Start(peerNode1_port)
	go peerNode2.Start(peerNode2_port)

	time.Sleep(1 * time.Second)

	if len(peerNode1.OpenConnections.Values) != 1 {
		t.Errorf("peerNode1 doesn't have peerNode2 in its openConnections set")
	}
	if len(peerNode2.OpenConnections.Values) != 1 {
		t.Errorf("peerNode2 doesn't have peerNode1 in its openConnections set")
	}
}









