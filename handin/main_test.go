package main

import (
	. "DNO/handin/Helper"
	"net"
	"reflect"
	"testing"
	"time"
)

/*
	debug_mode has to be set to false for tests to pass for some mystical reason
 */

func createPeerNode( shouldMockInput bool) PeerNode {
	return PeerNode{
		OpenConnections: 		SafeSet_Conn{   Values: make(map[net.Conn]bool) },
		PeersInArrivalOrder: 	Array_string{},
		MessagesSent:    		SafeSet_string{ Values: make(map[string  ]bool) },
		TestMock:        		Mock{ ShouldMockInput: shouldMockInput },
	}
}


func TestReceivedPeerListWhenJoining(t *testing.T) {

	const peerNode1_port = "50010"
	const peerNode2_port = "50011"
	const peerNode3_port = "50012"
	const peerNode4_port = "50013"

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)
	peerNode3 := createPeerNode(true)
	peerNode4 := createPeerNode(true)

	go peerNode1.Start(peerNode1_port, "no_port")
	time.Sleep(500 * time.Millisecond)
	go peerNode2.Start(peerNode2_port, peerNode1_port)
	time.Sleep(500 * time.Millisecond)
	go peerNode3.Start(peerNode3_port, peerNode1_port)
	time.Sleep(500 * time.Millisecond)
	go peerNode4.Start(peerNode4_port, peerNode1_port)
	time.Sleep(500 * time.Millisecond)

	expectedList := []string{peerNode1_port, peerNode2_port, peerNode3_port, peerNode4_port}
	if len(peerNode4.PeersInArrivalOrder.Values) == 0 {
		t.Errorf("peerNode4 didn't receive peerNode1's list of peers")
	}
	if !reflect.DeepEqual(peerNode4.PeersInArrivalOrder.Values, expectedList) {
		t.Errorf("peerNode4 received peerNode1's list of peers in the wrong order")
	}


}
func TestMessageIsSentFromNode1To3Via2(t *testing.T) {

	const peerNode1_port = "50007"
	const peerNode2_port = "50008"
	const peerNode3_port = "50009"

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)
	peerNode3 := createPeerNode(true)

	go peerNode1.Start(peerNode1_port, "no_port")
	time.Sleep(500 * time.Millisecond)
	go peerNode2.Start(peerNode2_port, peerNode1_port)
	time.Sleep(500 * time.Millisecond)
	go peerNode3.Start(peerNode3_port, peerNode2_port)
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

	go peerNode1.Start(peerNode1_port, "no_port")

	time.Sleep(1 * time.Second)
	peerNode1.TestMock.SimulatedInputString = "some_msg_1"
	time.Sleep(1 * time.Second)
	peerNode1.TestMock.SimulatedInputString = "some_msg_2"
	time.Sleep(1 * time.Second)
	peerNode1.TestMock.SimulatedInputString = "some_msg_3"
	time.Sleep(1 * time.Second)

	go peerNode2.Start(peerNode2_port, peerNode1_port)

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

	go peerNode1.Start(peerNode1_port, "no_port")
	go peerNode2.Start(peerNode2_port, peerNode1_port)

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

	go peerNode1.Start(peerNode1_port, "dont_connect_to_any_peer")
	go peerNode2.Start(peerNode2_port, peerNode1_port)

	time.Sleep(1 * time.Second)

	if len(peerNode1.OpenConnections.Values) != 1 {
		t.Errorf("peerNode1 doesn't have peerNode2 in its openConnections set")
	}
	if len(peerNode2.OpenConnections.Values) != 1 {
		t.Errorf("peerNode2 doesn't have peerNode1 in its openConnections set")
	}
}









