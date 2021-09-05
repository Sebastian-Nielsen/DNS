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
	OpenConnections SafeSet_Conn
	MessagesSent SafeSet_string
	Listener     net.Listener
	Ipc          IPC
	simulatedInputForTesting string
	TestMock                 Mock
}

var peerNode = PeerNode{
	OpenConnections:          SafeSet_Conn{   Values: make(map[net.Conn]bool) },
	MessagesSent:             SafeSet_string{ Values: make(map[string  ]bool) },
	Ipc:                      IPC{},
	simulatedInputForTesting: "",
	TestMock:                 Mock{ ShouldMockInput: false },
}

func (p* PeerNode) dialRemoteSocket() {
	// Dial the socket and add the connection to peerNode.OpenConnections if successful
	remoteSocket := PromptForRemoteSocket(p)

	conn, err := net.Dial("tcp", remoteSocket.ToString())

	dialIsSuccess := err == nil
	if !dialIsSuccess {
		fmt.Println("Dial failed")
		return
	} else {
		fmt.Println("--------------")
		fmt.Println("Our Ip:", conn.LocalAddr())
		fmt.Println("--------------")

		fmt.Println("Dial successful")
		peerNode.OpenConnections.Add(conn)
		go peerNode.Listen(conn)
	}
}
func (p* PeerNode) ListenForNewConns() {
	// Continously listen for new connection requests
	fmt.Println("[PeerNode] > Listening for new connections ...")
	for {
		fmt.Println("DEBUG 1")
		newConn, _ := p.Listener.Accept()
		fmt.Println("DEBUG 2")
		p.OpenConnections.Add(newConn)
		fmt.Println("[PeerNode] > Got a new connection ...")
		go p.Listen(newConn)
	}
}
func (p* PeerNode) Listen(conn net.Conn) {
	defer conn.Close()
	for {
		packet, ok := p.Ipc.Receive(conn)
		if !ok { return }
		go p.handleIncomming(packet, conn)
	}
}
func (p* PeerNode) handleIncomming(packet Packet, connPacketWasReceivedOn net.Conn) {

	switch packet.Type {
	case "UPDATE":
		if p.MessagesSent.Contains(packet.Msg) {
			fmt.Sprintf("Received msg we already have: %s", packet.Msg)
			return  // Ignore the packet
		}
		fmt.Println("Received packet: [Type: UPDATE][Msg: never_seen_before_msg] ... Broadcasting it")
		p.Broadcast(packet)
	case "PULL":
		fmt.Println("Received packet: [Type: PULL] ... Sending entire messagesSent-set back to sender")
		packet = Packet{ Type: "PULL-REPLY", MessagesSent: p.MessagesSent.Values }
		p.Ipc.Send(packet, connPacketWasReceivedOn)
	case "PULL-REPLY":
		fmt.Println("Received packet: [Type: PULL-REPLY] ... Extending our messagesSent set with the messages in the packet")
		p.HandlePullReplyPacket(packet)
	}
}
func (p * PeerNode) HandlePullReplyPacket(packet Packet) {
	// Add all messages contained in the PULL-REPLY packet to our set of messages
	for msg := range packet.MessagesSent {
		p.MessagesSent.Add(msg)
	}
}
func (p *PeerNode) HandleOutgoing(packet Packet) {

	switch packet.Type {
	case "UPDATE":
		if p.MessagesSent.Contains(packet.Msg) {
			fmt.Sprintf("Cancelling the sending of msg: '%s' (reason: already in messagesSent)", packet.Msg)
			return // Ignore packet
		}
		p.Broadcast(packet)
	default:
		fmt.Errorf("[HandleOutgoing] Got an unknown packet.Type: %s", packet.Type)
	}
}
func (p *PeerNode) Broadcast(packet Packet) {

	fmt.Println(p.MessagesSent.Values)
	p.MessagesSent.Add(packet.Msg)

	fmt.Println("Broadcasting msg: '" + packet.Msg + "'")
	for openConn := range p.OpenConnections.Values {
		fmt.Println("Sending msg to openConn:", openConn.RemoteAddr())
		p.Ipc.Send(packet, openConn)
	}
}
func (p *PeerNode) PullFromNeighbors() {
	// Requests the full messagesSent set from each neighbor in turn, wait inbetween each request
	packet := Packet{ Type: "PULL" }
	for {
		for openConn := range p.OpenConnections.Values {
			time.Sleep(30 * time.Second)
			fmt.Println("Sending pull-request to neighbor:", openConn.RemoteAddr().String())
			p.Ipc.Send( packet, openConn )
		}

	}
}

func (p *PeerNode) Start(port string) {
	p.Listener, _ = net.Listen("tcp", ":" + port)

	fmt.Println("[PeerNode:Start]", p.Listener.Addr())

	p.dialRemoteSocket() // Dial the socket and add the connection to peerNode.OpenConnections if successful

	go p.ListenForNewConns()

	go p.PullFromNeighbors() // Occassionally send pull-requests to neighbors, asking for their messagesSent set

	p.send() // Continously ask for user input to send
}

func main() {
	peerNode.Start("")


	fmt.Println("---------------------\nTerminating...")
}

func (p* PeerNode) send() {
	// Continously ask for input to send
	fmt.Println("-----------------------------------------")
	fmt.Println("[TcpClient] > Awaiting input to send ... ")
	fmt.Println("Type 'm' to view messagesSent ")
	fmt.Println("-----------------------------------------")
	for {
		msg := input(p)
		if strings.TrimSpace(msg) == "m" { fmt.Println(peerNode.MessagesSent.Values); continue }
		peerNode.HandleOutgoing( Packet{Type: "UPDATE", Msg: strings.TrimSpace(msg)} )
	}
}
func PromptForRemoteSocket(p* PeerNode) Socket {

	//fmt.Print("IP address:\n> ")
	//fmt.Scanln(&ip)
	//ip := input()
	ip := "127.0.0.1"  // The handin asks us to also prompt the user for an ip, but really no need for it ...


	fmt.Print("Port number: ")
	port := input(p)
	fmt.Printf("You wrote port: '%s'\n", port)


	return Socket{ip, port}
}

func input(p* PeerNode) string {
	if !p.TestMock.ShouldMockInput {
		var returnString string
		fmt.Scanln(&returnString)
		return returnString
	}

	for p.TestMock.SimulatedInputString == "" {} // waiting for input
	returnString := p.TestMock.SimulatedInputString
	p.TestMock.SimulatedInputString = ""
	return returnString

}
