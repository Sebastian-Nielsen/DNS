package main

import (
	. "DNO/handin/Account"
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
	OpenConnections     SafeSet_Conn
	PeersInArrivalOrder SafeArray_string
	MessagesSent        SafeSet_string
	Listener            net.Listener
	Ipc                 IPC
	TestMock            Mock
	LocalLedger         *Ledger
	TransactionsSeen   SafeArray_Transaction
}

func (p *PeerNode) createSocket(remotePort string) Socket {
	if remotePort != "" {
		return Socket{Ip: "127.0.0.1", Port: remotePort}
	}
	return PromptForRemoteSocket(p)
}
func (p *PeerNode) DialNetwork(remoteSocket Socket) {
	conn, err := net.Dial("tcp", remoteSocket.ToString())

	dialIsSuccess := err == nil
	if !dialIsSuccess {
		p.debugPrintln("--------------")
		p.println("DialNetwork failed")
		p.debugPrintln("--------------")
		p.debugPrintln("Adding local port to list")
		p.PeersInArrivalOrder.Append(PortOf(p.Listener.Addr()))
		return
	} else {
		p.debugPrintln("--------------")
		p.debugPrintln("Local  Addr:", conn.LocalAddr())
		p.debugPrintln("Remote Addr:", conn.RemoteAddr())
		p.debugPrintln("--------------")
		p.OpenConnections.Add(conn)

		p.debugPrintln("[dial] Sending Listener port to dialed connection (" + PortOf(conn.RemoteAddr()) + ")")

		//establishConnectionToNetwork()
		//setup

		p.Ipc.Send(
			Packet{Type: PacketType.LISTENER_PORT, Msg: PortOf(p.Listener.Addr()) },
			conn,
		)

		p.debugPrintln("Sending initial pull-request to neighbor:", conn.RemoteAddr().String())

		p.Ipc.Send(Packet{Type: PacketType.PULL}, conn)
		var packet Packet
		for {
			packet, _ = p.Ipc.Receive(conn)
			if packet.Type == PacketType.PULL_REPLY {break}
			p.debugPrintf("[Dial] Received wrong packet type:", packet.Type)
		}
		p.debugPrintln("[Dial] Received Pull-Reply packet:", packet.PeersInArrivalOrderValues, packet.MessagesSent)

		p.HandlePullReplyPacket(packet)

		peers := packet.PeersInArrivalOrderValues

		p.connectToPeers(peers)

		p.Broadcast(
			Packet{ Type: PacketType.BROADCAST_LISTENER_PORT , Msg: PortOf(p.Listener.Addr()) },
		)

		go p.Listen(conn)}
}
func Assert(condition bool) {
	if !condition {
		panic("assert failed")
	}
}
func (p *PeerNode) connectToPeers(portsOfPeers []string) {
	Assert(p.Listener != nil)
	//upTo := int(math.Min(11, float64(len(portsOfPeers)))-1)

	p.debugPrintln("Connecting to up to 10 peers in received list:\n\t", portsOfPeers)

	if len(portsOfPeers) == 0 {
		p.println("Error: Received empty list from connection")
		return
	}
	newConnections := 0
	for _, port := range portsOfPeers {       // portsOfPeers = ["50001", "50002", ...]
		alreadyHaveAConnectionToThisPort := p.OpenConnections.ContainsAConnWith(port)
		if alreadyHaveAConnectionToThisPort {
			p.debugPrintln("Don't connect to this, we already have a connection to port:", port)
			continue
		}
		isPortOurRemotePort := port == PortOf(p.Listener.Addr())
		if isPortOurRemotePort {
			p.debugPrintln("Don't connect to this, it is ourselves *insert spiderman pointing meme*:", port)
			continue
		}

		conn := p.dial(Socket{Ip: "127.0.0.1", Port: port})
		p.OpenConnections.Add(conn)
		newConnections += 1
		p.debugPrint("\tConnecting to peer:", port)
		if newConnections == 10 {
			break
		}
	}
}
func (p *PeerNode) dial(socket Socket) net.Conn {
	conn, err := net.Dial("tcp", socket.ToString())
	if err != nil {
		p.println(err)
	}
	return conn
}
func (p *PeerNode) ListenForNewConns() {
	// Continously listen for new connection requests
	p.debugPrintln("Listening for new connections ...")
	for {
		newConn, _ := p.Listener.Accept()
		p.OpenConnections.Add(newConn)
		p.debugPrintln("Got a new connection ...", newConn.LocalAddr(), "->", newConn.RemoteAddr())
		go p.Listen(newConn)
	}
}
func (p *PeerNode) Listen(conn net.Conn) {
	defer conn.Close()
	for {
		packet, ok := p.Ipc.Receive(conn)
		if !ok { p.debugPrintln("Receive error") }
		go p.handleIncomming(packet, conn)
	}
}
func (p *PeerNode) handleIncomming(packet Packet, connPacketWasReceivedOn net.Conn) {

	switch packet.Type {
	case PacketType.BROADCAST_MSG:
		p.debugPrintln("Received packet: [Type: BROADCAST_MSG][Msg: never_seen_before_msg] ... Broadcasting it")
		p.BroadcastMessage(packet)
	case PacketType.PULL:
		packet = Packet{
			Type: PacketType.PULL_REPLY,
			MessagesSent: p.MessagesSent.Values,
			PeersInArrivalOrderValues: p.PeersInArrivalOrder.Values(),
			TransactionsSeen: p.TransactionsSeen.Values(),
		}
		p.debugPrintln("Received packet: [Type: PULL] ... Sending back: \n\tpacket{peersInArrivalOrder:", packet.PeersInArrivalOrderValues, "}")
		p.Ipc.Send(packet, connPacketWasReceivedOn)
	case PacketType.PULL_REPLY:
		p.debugPrintln("Received packet: [Type: PULL-REPLY] ... Extending our messagesSent set with the messages in the packet")
		p.HandlePullReplyPacket(packet)
	case PacketType.LISTENER_PORT:
		listenerPort := packet.Msg
		p.PeersInArrivalOrder.Append(listenerPort)
		p.debugPrintln("Received packet: [Type: LISTENER_PORT] ... Adding its port:", listenerPort, "peersInArrivalOrder is now:\n\t", p.PeersInArrivalOrder.Values())
	case PacketType.BROADCAST_LISTENER_PORT:
		listenerPort := packet.Msg
		p.debugPrintln("received packet: [Type: BROADCAST_LISTENER_PORT] ... ", listenerPort)
		if !p.PeersInArrivalOrder.Contains(listenerPort) {
			p.PeersInArrivalOrder.Append(listenerPort)
			p.debugPrintln("[PacketType: BROADCASTED_LISTENER_PORT] Adding port:", listenerPort, "peersInArrivalOrder is now:\n\t", p.PeersInArrivalOrder.Values())
		} else {
			packet.Type = PacketType.BROADCASTED_KNOWN_LISTENER_PORT
		}
		p.Broadcast(packet)
	case PacketType.BROADCASTED_KNOWN_LISTENER_PORT:
		listenerPort := packet.Msg
		p.debugPrintln("Received packet: [Type: BROADCASTED_KNOWN_LISTENER_PORT] ... ", listenerPort)

		if p.PeersInArrivalOrder.Contains(listenerPort) {
			p.debugPrintf("Received listenerPort we already have: %s\n", packet.Msg)
			return // Ignore the packet
		}
		p.PeersInArrivalOrder.Append(listenerPort)
		p.debugPrintln("[PacketType: BROADCASTED_KNOWN_LISTENER_PORT] Adding port", listenerPort, "peersInArrivalOrder is now:\n\t", p.PeersInArrivalOrder.Values())
		p.Broadcast(packet)
	case PacketType.BROADCAST_TRANSACTION:
		p.debugPrintln("received packet: [PacketType: BROADCASTED_TRANSACTION]", packet.Transaction)
		transaction := packet.Transaction

		transactionIsNotSeen := !p.TransactionsSeen.Contains(transaction)
		if transactionIsNotSeen {
			p.TransactionsSeen.Append(transaction)
			p.LocalLedger.Transaction(&transaction)
		}

		packet.Type = PacketType.BROADCAST_KNOWN_TRANSACTION
		p.Broadcast(packet)
	case PacketType.BROADCAST_KNOWN_TRANSACTION:
		p.debugPrintln("[PacketType: BROADCASTED__KNOWN_TRANSACTION]", packet.Transaction)
		transaction := packet.Transaction

		transactionIsNotSeen := !p.TransactionsSeen.Contains(transaction)
		if transactionIsNotSeen {
			p.TransactionsSeen.Append(transaction)
			p.LocalLedger.Transaction(&transaction)
			p.Broadcast(packet)
		}
	}
}
func (p *PeerNode) HandlePullReplyPacket(packet Packet) {
	// Add all messages contained in the PULL-REPLY packet to our set of messages
	p.extendMessagesSentSet(packet.MessagesSent)

	// Add all peer ports contained in the PULL-REPLY packet to our list of ports
	p.extendPeersInArrivalOrder(packet.PeersInArrivalOrderValues)

	// Apply all transactions contained in the PULL-REPLY packet on our local ledger
	p.applyAllTransactions(packet.TransactionsSeen)
}
func (p *PeerNode) applyAllTransactions(transactions []Transaction) {
	p.debugPrintln("Applying all transactions:", transactions)
	for _, transaction := range transactions {
		p.LocalLedger.Transaction(&transaction)
		p.TransactionsSeen.Append(transaction)
	}
}
func (p *PeerNode) extendPeersInArrivalOrder(peers []string) {
	// no contains method i golang, so simple O(n*m) solution
	p.debugPrintln("Extending PeersInArrivalOrder")
	for _, receivedPeer := range peers {
		if !p.PeersInArrivalOrder.Contains(receivedPeer) {
			p.debugPrint("\tAdding:", receivedPeer, "to PeersInArrivalOrder")
			p.PeersInArrivalOrder.Append(receivedPeer)
		}
	}
	p.debugPrintln("PeersInArrivalOrder is now:", p.PeersInArrivalOrder.Values())
}
func (p *PeerNode) extendMessagesSentSet(messages map[string]bool) {
	if len(messages) == 0 { return }

	p.debugPrintln("Adding messages:")
	for msg := range messages {
		p.println("\t", msg)
		p.MessagesSent.Add(msg)
	}
}
func (p *PeerNode) HandleOutgoing(packet Packet) {
	switch packet.Type {
	case PacketType.BROADCAST_MSG:
		if p.MessagesSent.Contains(packet.Msg) {
			p.debugPrintf("Cancelling the sending of msg: '%s' (reason: already in messagesSent)", packet.Msg)
			return // Ignore packet
		}
		p.MessagesSent.Add(packet.Msg)
		p.Broadcast(packet)
	default:
		fmt.Printf("[HandleOutgoing] Got an unknown packet.Type: %s", packet.Type)
	}
}
func (p *PeerNode) BroadcastMessage(packet Packet) {
	if p.MessagesSent.Contains(packet.Msg) {
		p.debugPrintf("Received msg we already have: %s", packet.Msg)
		return // Ignore the packet
	}
	p.println("Adding and broadcasting msg: '" + packet.Msg + "'")
	p.MessagesSent.Add(packet.Msg)
	p.Broadcast(packet)
}
func (p *PeerNode) Broadcast(packet Packet) {
	p.debugPrintf("[There's now %d openConnections:\n\t%s\n", len(p.OpenConnections.Values), p.OpenConnections.ToString())
	for openConn := range p.OpenConnections.Values {
		p.println("Sending msg to openConn:", openConn.RemoteAddr())
		p.Ipc.Send(packet, openConn)
	}
}
func (p *PeerNode) PullFromNeighbors() {
	// Requests the full messagesSent set from each neighbor in turn, wait inbetween each request
	time.Sleep(300 * time.Second)
	packet := Packet{Type: PacketType.PULL}
	for {
		for openConn := range p.OpenConnections.Values {
			p.println("Sending pull-request to neighbor:", openConn.RemoteAddr().String())
			p.Ipc.Send(packet, openConn)
			time.Sleep(30 * time.Second)
		}
	}
}
func (p *PeerNode) startServer(port string) {
	p.Listener, _ = net.Listen("tcp", ":" + port)
	p.printPort(p.Listener)

	go p.ListenForNewConns()
}
func (p *PeerNode) send() {
	/* Continously prompt the user for messages to send */
	//p.debugPrintln("[PeerNode:send] Awaiting input to send ... ")
	//p.println("[PeerNode:send] > Type 'm' to view MessagesSent ")
	//p.println("[PeerNode:send] > Type 'c' to view OpenConnections ")
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
		if msg == "a" {
			println(p.PeersInArrivalOrder.Values())
			continue
		}
		p.HandleOutgoing(Packet{Type: PacketType.BROADCAST_MSG, Msg: msg})
	}
}


/*
	Transaction related methods
 */
func (p *PeerNode) CreateAccountInLedger(id string, initialAmount int) {
	//p.LocalLedger.CreateAccount(id, initialAmount)
	p.LocalLedger.Accounts[id] = initialAmount
}
func (p *PeerNode) MakeAndBroadcastTransaction(amount int, id string, from string, to string) {
	transaction := Transaction{ID: id, From: from, To: to, Amount: amount}

	p.debugPrintln("Applying transaction:", transaction, ". Accounts before:", p.LocalLedger.Accounts)
	p.LocalLedger.Transaction(&transaction)
	p.TransactionsSeen.Append(transaction)
	p.debugPrintln("Applying transaction:", transaction, ". Accounts after:", p.LocalLedger.Accounts)

	p.Broadcast(
		Packet{
			Type: PacketType.BROADCAST_TRANSACTION,
			Transaction: transaction,
		},
	)
}

/*
	atPort: port at which the PeerNode should listen at.
	remotePort: port of remote PeerNode to connect to initially.
				if remotePort is "" (empty string) then prompt the
				user for a port.
 */
func (p *PeerNode) Start(atPort, remotePort string) {
	p.startServer(atPort)

	p.DialNetwork(p.createSocket(remotePort))

	go p.PullFromNeighbors() // Occassionally send pull-requests to neighbors, asking for their messagesSent set

	p.send() // Continously prompt the user msgs for the peerNode to broadcast
}

const DEBUG_MODE = true
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

func PromptForRemoteSocket(p *PeerNode) Socket {

	ip := "127.0.0.1"  // The handin asks us to also prompt the user for an ip, but really no need for it ...

	fmt.Print("Connect to port: ")
	port := strings.TrimSpace(input(p))
	p.debugPrintf("You wrote port: '%s'\n", port)

	return Socket{Ip: ip, Port: port}
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



/*
	Printer helper methods
 */
func (p *PeerNode) debugPrintf(text string, args ...interface{}) {
	if p.TestMock.ShouldPrintDebug {
		fmt.Printf("<" + PortOf(p.Listener.Addr()) + "> " + text, args...)
	}
}
func (p *PeerNode) debugPrint(args ...interface{}) {
	if p.TestMock.ShouldPrintDebug {
		fmt.Print( "\t", args, "\n")
	}
}
func (p *PeerNode) debugPrintln(args ...interface{}) {
	if p.TestMock.ShouldPrintDebug {
		fmt.Println("<" + PortOf(p.Listener.Addr()) + ">", args)
	}
}
func (p *PeerNode) println(args ...interface{}) {
	fmt.Println("<" + PortOf(p.Listener.Addr()) + ">", args)
}
func (p *PeerNode) printPort(listener net.Listener) {
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	p.println("Running on port: " + port)
}

