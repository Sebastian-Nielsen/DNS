package main

import (
	. "DNO/handin/Helper"
	"fmt"
	"net"
	"strings"
	"time"
)

/*
	PeerNode
	- Contionusly listen for new connection requests
		- add new connections to OpenConnections

	- process incomming messages
	- process outgoing messages

	- Occassionally send a pull request to a neighbor (taking turns)
*/
type PeerNode struct {
	OpenConnections          SafeSet_Conn
	MessagesSent             SafeSet_string
	Listener                 net.Listener
	Ipc                      IPC
	simulatedInputForTesting string
	TestMock                 Mock
}

func (p *PeerNode) dialRemoteSocket() {
	// prompt the user for a socket to dial; add the connection to peerNode.OpenConnections if successful
	remoteSocket := PromptForRemoteSocket(p)

	conn, err := net.Dial("tcp", remoteSocket.ToString())

	dialIsSuccess := err == nil
	if !dialIsSuccess {
		fmt.Println("Dial failed")
		return
	} else {
		println("--------------")
		println("Our Ip:", conn.LocalAddr())
		println("--------------")

		fmt.Println("Dial successful")
		p.OpenConnections.Add(conn)
		go p.Listen(conn)
	}
}
func (p *PeerNode) ListenForNewConns() {
	// Continously listen for new connection requests
	println("Listening for new connections ...")
	for {
		newConn, _ := p.Listener.Accept()
		p.OpenConnections.Add(newConn)
		println("Got a new connection ...")
		go p.Listen(newConn)
	}
}
func (p *PeerNode) Listen(conn net.Conn) {
	defer conn.Close()
	for {
		packet, ok := p.Ipc.Receive(conn)
		if !ok { return }
		go p.handleIncomming(packet, conn)
	}
}
func (p *PeerNode) handleIncomming(packet Packet, connPacketWasReceivedOn net.Conn) {

	switch packet.Type {
	case "UPDATE":
		if p.MessagesSent.Contains(packet.Msg) {
			printf("Received msg we already have: %s", packet.Msg)
			return // Ignore the packet
		}
		println("Received packet: [Type: UPDATE][Msg: never_seen_before_msg] ... Broadcasting it")
		p.Broadcast(packet)
	case "PULL":
		println("Received packet: [Type: PULL] ... Sending entire messagesSent-set back to sender")
		packet = Packet{Type: "PULL-REPLY", MessagesSent: p.MessagesSent.Values}
		p.Ipc.Send(packet, connPacketWasReceivedOn)
	case "PULL-REPLY":
		println("Received packet: [Type: PULL-REPLY] ... Extending our messagesSent set with the messages in the packet")
		p.HandlePullReplyPacket(packet)
	}
}
func (p *PeerNode) HandlePullReplyPacket(packet Packet) {
	// Add all messages contained in the PULL-REPLY packet to our set of messages
	if len(packet.MessagesSent) != 0 {
		fmt.Println("Adding messages:")
	}
	for msg := range packet.MessagesSent {
		fmt.Printf("\t'%s'", msg)
		p.MessagesSent.Add(msg)
	}

}
func (p *PeerNode) HandleOutgoing(packet Packet) {

	switch packet.Type {
	case "UPDATE":
		if p.MessagesSent.Contains(packet.Msg) {
			printf("Cancelling the sending of msg: '%s' (reason: already in messagesSent)", packet.Msg)
			return // Ignore packet
		}
		p.Broadcast(packet)
	default:
		fmt.Errorf("[HandleOutgoing] Got an unknown packet.Type: %s", packet.Type)
	}
}
func (p *PeerNode) Broadcast(packet Packet) {
	fmt.Printf("Adding msg: '%s'\n", packet.Msg)
	p.MessagesSent.Add(packet.Msg)

	fmt.Printf("Broadcasting msg: '%s'\n", packet.Msg)
	printf("[Broadcasting] to %d clients\n", len(p.OpenConnections.Values))
	for openConn := range p.OpenConnections.Values {
		println("Sending msg to openConn:", openConn.RemoteAddr())
		p.Ipc.Send(packet, openConn)
	}
}
func (p *PeerNode) PullFromNeighbors() {
	// Requests the full messagesSent set from each neighbor in turn, wait inbetween each request
	packet := Packet{Type: "PULL"}
	for {
		for openConn := range p.OpenConnections.Values {
			println("Sending pull-request to neighbor:", openConn.RemoteAddr().String())
			p.Ipc.Send(packet, openConn)
			time.Sleep(30 * time.Second)
		}

	}
}
func (p *PeerNode) startServer(port string) {
	p.Listener, _ = net.Listen("tcp", ":" + port)
	printPort(p.Listener)
	go p.ListenForNewConns()
}
func (p *PeerNode) Start(atPort string) {
	p.dialRemoteSocket() // prompt the user for a socket to dial; add the connection to peerNode.OpenConnections if successful

	go p.startServer(atPort) // Starts a connection and then starts listening for new connections

	go p.PullFromNeighbors() // Occassionally send pull-requests to neighbors, asking for their messagesSent set

	p.send() // Continously prompt the user msgs for the peerNode to broadcast
}

const DEBUG_MODE = false
func main() {
	var peerNode = PeerNode{
		OpenConnections:          SafeSet_Conn{   Values: make(map[net.Conn]bool) },
		MessagesSent:             SafeSet_string{ Values: make(map[string  ]bool) },
		Ipc:                      IPC{},
		simulatedInputForTesting: "",
		TestMock:                 Mock{ ShouldMockInput: false },
	}

	peerNode.Start("")

	println("---------------------\nTerminating...")
}

func (p* PeerNode) send() {
	/* Continously prompt the user for messages to send */
	fmt.Println("[PeerNode:send] Awaiting input to send ... ")
	//fmt.Println("[PeerNode:send] > Type 'm' to view MessagesSent ")
	//fmt.Println("[PeerNode:send] > Type 'c' to view OpenConnections ")
	for {
		msg := strings.TrimSpace(input(p))
		if msg == "m" {
			println(p.MessagesSent.Values)
			continue
		}
		if msg == "c" {
			println(p.OpenConnections.Values)
			continue
		}
		p.HandleOutgoing(Packet{Type: "UPDATE", Msg: msg})
	}
}
func PromptForRemoteSocket(p *PeerNode) Socket {

	ip := "127.0.0.1"  // The handin asks us to also prompt the user for an ip, but really no need for it ...

	fmt.Print("Connect to port: ")
	port := strings.TrimSpace(input(p))
	printf("You wrote port: '%s'\n", port)

	return Socket{ip, port}
}
func input(p *PeerNode) string {
	if !p.TestMock.ShouldMockInput {
		var returnString string
		fmt.Scanln(&returnString)
		return returnString
	}

	for p.TestMock.SimulatedInputString == "" {
		time.Sleep(250 * time.Millisecond)
	} // waiting for input
	returnString := p.TestMock.SimulatedInputString
	p.TestMock.SimulatedInputString = ""
	return returnString

}

func printf(text string, args ...interface{}) {
	if DEBUG_MODE {
		fmt.Printf(text, args...)
	}
}
func println(args ...interface{}) {
	if DEBUG_MODE {
		fmt.Println(args...)
	}
}

func printPort(listener net.Listener) {
	connAddr := listener.Addr().String()
	portIndex := strings.LastIndex(connAddr, ":")
	fmt.Println("Running on port: " + connAddr[portIndex+1:])
}