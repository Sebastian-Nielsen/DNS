package main

import (
	. "DNO/handin/Helper"
	"net"
	"reflect"
	"testing"
	"time"
)

/*
	sometimes debug_mode has to be set to false for tests to pass for some mystical reason
 */

func createPeerNode( shouldMockInput bool) PeerNode {
	return PeerNode{
		OpenConnections:     SafeSet_Conn{   Values: make(map[net.Conn]bool) },
		PeersInArrivalOrder: SafeArray_string{},
		MessagesSent:        SafeSet_string{ Values: make(map[string  ]bool) },
		TestMock:            Mock{ ShouldMockInput: shouldMockInput },
	}
}


func TestNodeConnectsToThreeOthersWhenEnteringNetwork(t *testing.T) {
	t.Parallel()

	peerNode1 := createPeerNode(true)
	peer1Port := AvailablePorts.Next()
	goStart(&peerNode1, peer1Port, "no_port")

	var peerNode2 = createPeerNode(true)
	var peerNode3 = createPeerNode(true)
	var peerNode4 = createPeerNode(true)

	goStart(&peerNode2, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode3, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode4, AvailablePorts.Next(), peer1Port)

	time.Sleep(10 * time.Second)

	if !reflect.DeepEqual(peerNode4.OpenConnections.Values, peerNode1.OpenConnections.Values) {
		t.Error("peerNode1 and peerNode4 has different connections.\n" +
			"PeerNode1's openConnections:", peerNode1.OpenConnections.ToString(), "\n",
			"\nPeerNode4's openConnections:", peerNode4.OpenConnections.ToString())
	}

	//p12Port := AvailablePorts.Next()
	//peerNode12 := createPeerNode(true)
	//goStart(&peerNode12, p12Port, peer1Port)

	//time.Sleep(5 * time.Second)
	//
	//if !peerNode4.PeersInArrivalOrder.Contains(p4Port) {
	//	t.Error("peerNode3 doesn't have the port of peerNode4 (" + p4Port + ") after peerNode1 broadcast it.\n" +
	//		"PeerNode3's portList:", peerNode3.PeersInArrivalOrder.Values())
	//}

	time.Sleep(2 * time.Second)

}

func TestNodeConnectsToTenOthersWhenEnteringNetwork(t *testing.T) {
	t.Parallel()

	peerNode1 := createPeerNode(true)
	peer1Port := AvailablePorts.Next()
	goStart(&peerNode1, peer1Port, "no_port")

	var peerNode2 = createPeerNode(true)
	var peerNode3 = createPeerNode(true)
	var peerNode4 = createPeerNode(true)
	var peerNode5 = createPeerNode(true)
	var peerNode6 = createPeerNode(true)
	var peerNode7 = createPeerNode(true)
	var peerNode8 = createPeerNode(true)
	var peerNode9 = createPeerNode(true)
	var peerNode10 = createPeerNode(true)
	var peerNode11 = createPeerNode(true)

	goStart(&peerNode2, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode3, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode4, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode5, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode6, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode7, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode8, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode9, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode10, AvailablePorts.Next(), peer1Port)
	goStart(&peerNode11, AvailablePorts.Next(), peer1Port)

	//var currentPeerNode PeerNode
	//for range [10]int{} {  // Create 10 peerNodes and connect them all to peerNode1
	//	currentPeerNode = createPeerNode(true)
	//	currentPeerPort := AvailablePorts.Next()
	//	goStart(&currentPeerNode, currentPeerPort, peer1Port)
	//}
	//var peerNode11 = currentPeerNode

	time.Sleep(20 * time.Second)

	if !reflect.DeepEqual(peerNode11.OpenConnections.Values, peerNode1.OpenConnections.Values) {
		t.Error("peerNode1 and peerNode11 has different connections.\n" +
			"PeerNode1's openConnections:", peerNode1.OpenConnections.ToString(), "\n",
			"\nPeerNode11's openConnections:", peerNode11.OpenConnections.ToString())
	}

	p12Port := AvailablePorts.Next()
	peerNode12 := createPeerNode(true)
	goStart(&peerNode12, p12Port, peer1Port)

	time.Sleep(5 * time.Second)

	if !peerNode11.PeersInArrivalOrder.Contains(p12Port) {
		t.Error("peerNode11 doesn't have the port of peerNode12 (" + p12Port + ") after peerNode1 broadcast it.\n" +
			"PeerNode11's portList:", peerNode11.PeersInArrivalOrder.Values())
	}

	time.Sleep(2 * time.Second)

}


func TestPeerNodeConnectsToAllNodesWhenEnteringNetwork(t *testing.T) {
	//t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()
	var peerNode4_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)
	peerNode3 := createPeerNode(true)
	peerNode4 := createPeerNode(true)

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode1_port)
	goStart(&peerNode4, peerNode4_port, peerNode1_port)

	time.Sleep(2 * time.Second)
	p4Conns := peerNode4.OpenConnections.Values
	if len(p4Conns) == 3 {
		t.Error("peerNode4 doesn't have an open connection to each node in the network.\n" +
			"PeerNode4's openConnections:", p4Conns)
	}
}
func TestReceivedPeerListWhenJoining(t *testing.T) {
	//t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()
	var peerNode4_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)
	peerNode3 := createPeerNode(true)
	peerNode4 := createPeerNode(true)

	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode1_port)
	goStart(&peerNode4, peerNode4_port, peerNode1_port)

	time.Sleep(2 * time.Second)

	expectedList := []string{peerNode1_port, peerNode2_port, peerNode3_port, peerNode4_port}
	if len(peerNode4.PeersInArrivalOrder.Values()) == 0 {
		t.Errorf("peerNode4 didn't receive peerNode1's list of peers")
	}
	if !reflect.DeepEqual(peerNode4.PeersInArrivalOrder.Values(), expectedList) {
		t.Error("peerNode4 received peerNode1's list of peers in the wrong order:\npeerNode4:",
			peerNode4.PeersInArrivalOrder.Values(), "\nvs\npeerNode1 (expected list):", expectedList)
	}


}
func TestMessageIsSentFromNode1To3Via2(t *testing.T) {
	//t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()
	var peerNode3_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)
	peerNode3 := createPeerNode(true)


	goStart(&peerNode1, peerNode1_port, "no_port")
	goStart(&peerNode2, peerNode2_port, peerNode1_port)
	goStart(&peerNode3, peerNode3_port, peerNode2_port)


	simulateInputFor(&peerNode1, "some_msg")

	if !peerNode3.MessagesSent.Contains("some_msg") {
		t.Errorf("peerNode3 didn't receive peerNode1's msgs")
	}
}
func TestLatercomerNodeEventuallyGetsAllMsgs(t *testing.T) {
	//t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)

	go peerNode1.Start(peerNode1_port, "no_port")

	time.Sleep(1 * time.Second)

	simulateInputFor(&peerNode1, "some_msg_1")
	simulateInputFor(&peerNode1, "some_msg_2")
	simulateInputFor(&peerNode1, "some_msg_3")

	go peerNode2.Start(peerNode2_port, peerNode1_port)

	time.Sleep(3 * time.Second)

	if  !peerNode2.MessagesSent.Contains("some_msg_1") ||
		!peerNode2.MessagesSent.Contains("some_msg_2") ||
		!peerNode2.MessagesSent.Contains("some_msg_3") {
		t.Errorf("peerNode2 didn't receive peerNode1's msgs")
	}
}
func TestPeer1ReceivesMsgFromPeer2(t *testing.T) {
	//t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

	peerNode1 := createPeerNode(true)
	peerNode2 := createPeerNode(true)

	go peerNode1.Start(peerNode1_port, "no_port")
	go peerNode2.Start(peerNode2_port, peerNode1_port)

	time.Sleep(1 * time.Second)

	simulateInputFor(&peerNode1, "some_msg")

	if !peerNode2.MessagesSent.Contains("some_msg") {
		t.Errorf("peerNode2 didn't receive peerNode1's msg")
	}

}
func TestPeer1CanConnectToPeer2(t *testing.T) {
	//t.Parallel()

	var peerNode1_port = AvailablePorts.Next()
	var peerNode2_port = AvailablePorts.Next()

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



func simulateInputFor(node *PeerNode, text string) {
	node.TestMock.SimulatedInputString = text
	time.Sleep(2 * time.Second)
}
func goStart(peerNode *PeerNode, atPort string, remotePort string) {
	go peerNode.Start(atPort, remotePort)
	time.Sleep(20 * time.Second)
}


